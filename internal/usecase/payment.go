package usecase

import (
	"errors"
	"log"
	"shop/internal/domain"
	"shop/internal/repository"
	"time"
)

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
	userRepo    repository.UserRepository
	cartRepo    repository.CartItemsRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewPaymentUseCase(paymentRepo repository.PaymentRepository, cartRepo repository.CartItemsRepository, orderRepo repository.OrderRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		cartRepo:    cartRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// InitiatePayment - شروع فرآیند پرداخت و قفل کردن سبد خرید
func (p *PaymentUsecase) InitiatePayment(userID uint, paymentMethod string) (*domain.Payment, error) {
	// دریافت کاربر
	user, err := p.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// بررسی اینکه سبد خرید قفل نباشد
	if user.LockCart {
		return nil, errors.New("cart is already locked for payment")
	}

	// دریافت آیتم‌های سبد خرید
	cartItems, err := p.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// محاسبه مجموع قیمت و بررسی موجودی
	var total float64
	for _, item := range cartItems {
		product, err := p.productRepo.GetByID(item.ProductId)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		total += product.Price * float64(item.Quantity)
	}

	// قفل کردن سبد خرید
	user.LockCart = true
	user.TotalCart = total
	if err := p.userRepo.Update(user); err != nil {
		return nil, err
	}

	// ایجاد رکورد پرداخت
	payment := &domain.Payment{
		UserID:        userID,
		Amount:        total,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := p.paymentRepo.Create(payment); err != nil {
		// در صورت خطا، قفل سبد خرید را باز کن
		user.LockCart = false
		p.userRepo.Update(user)
		return nil, err
	}

	return payment, nil
}

// ProcessPayment - پردازش پرداخت (شبیه‌سازی gateway)
func (p *PaymentUsecase) ProcessPayment(paymentID uint, gatewayResponse domain.PaymentGatewayResponse) error {
	payment, err := p.paymentRepo.GetByID(paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != "pending" {
		return errors.New("payment is not in pending state")
	}

	// شبیه‌سازی پردازش در gateway
	if gatewayResponse.Success {
		payment.Status = "completed"
		payment.GatewayTransactionID = gatewayResponse.TransactionID
		payment.GatewayResponse = gatewayResponse.Message
	} else {
		payment.Status = "failed"
		payment.GatewayResponse = gatewayResponse.Message
	}

	payment.UpdatedAt = time.Now()

	if err := p.paymentRepo.Update(payment); err != nil {
		return err
	}

	// اگر پرداخت موفق بود، order ایجاد کن
	if payment.Status == "completed" {
		if err := p.createOrderFromCart(payment.UserID); err != nil {
			log.Printf("Error creating order after successful payment: %v", err)
			// در اینجا می‌توانید payment را به حالت failed تغییر دهید یا retry logic اضافه کنید
		}
	} else {
		// اگر پرداخت ناموفق بود، قفل سبد خرید را باز کن
		if err := p.unlockCart(payment.UserID); err != nil {
			log.Printf("Error unlocking cart after failed payment: %v", err)
		}
	}

	return nil
}

// createOrderFromCart - ایجاد سفارش از آیتم‌های سبد خرید
func (p *PaymentUsecase) createOrderFromCart(userID uint) error {
	// دریافت آیتم‌های سبد خرید
	cartItems, err := p.cartRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	if len(cartItems) == 0 {
		return errors.New("cart is empty")
	}

	// محاسبه مجموع و ایجاد order items
	var total float64
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		product, err := p.productRepo.GetByID(item.ProductId)
		if err != nil {
			return err
		}

		// بررسی نهایی موجودی
		if product.Stock < item.Quantity {
			return errors.New("insufficient stock")
		}

		orderItem := domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		total += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, orderItem)
	}

	// ایجاد سفارش
	order := &domain.Order{
		UserID:    userID,
		Total:     total,
		Status:    "confirmed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := p.orderRepo.Create(order); err != nil {
		return err
	}

	// ایجاد order items و به‌روزرسانی موجودی
	for _, item := range orderItems {
		item.OrderID = order.ID
		if err := p.orderRepo.CreateOrderItem(&item); err != nil {
			return err
		}

		// کاهش موجودی محصول
		product, _ := p.productRepo.GetByID(item.ProductID)
		newStock := product.Stock - item.Quantity
		if err := p.productRepo.UpdateStock(item.ProductID, newStock); err != nil {
			return err
		}
	}

	// پاک کردن سبد خرید و باز کردن قفل
	if err := p.clearCartAndUnlock(userID); err != nil {
		return err
	}

	return nil
}

// clearCartAndUnlock - پاک کردن سبد خرید و باز کردن قفل
func (p *PaymentUsecase) clearCartAndUnlock(userID uint) error {
	// پاک کردن آیتم‌های سبد خرید
	if err := p.cartRepo.ClearCart(userID); err != nil {
		return err
	}

	// باز کردن قفل سبد خرید و صفر کردن مجموع
	user, err := p.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.LockCart = false
	user.TotalCart = 0

	return p.userRepo.Update(user)
}

// unlockCart - فقط باز کردن قفل سبد خرید (برای حالت ناموفق)
func (p *PaymentUsecase) unlockCart(userID uint) error {
	user, err := p.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.LockCart = false
	return p.userRepo.Update(user)
}

// GetPayment - دریافت اطلاعات پرداخت
func (p *PaymentUsecase) GetPayment(paymentID uint) (*domain.Payment, error) {
	return p.paymentRepo.GetByID(paymentID)
}

// GetUserPayments - دریافت تاریخچه پرداخت‌های کاربر
func (p *PaymentUsecase) GetUserPayments(userID uint, page, limit int) ([]*domain.Payment, error) {
	offset := (page - 1) * limit
	return p.paymentRepo.GetByUserID(userID, limit, offset)
}

// CancelPayment - لغو پرداخت و باز کردن قفل سبد
func (p *PaymentUsecase) CancelPayment(paymentID uint) error {
	payment, err := p.paymentRepo.GetByID(paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != "pending" {
		return errors.New("only pending payments can be cancelled")
	}

	payment.Status = "cancelled"
	payment.UpdatedAt = time.Now()

	if err := p.paymentRepo.Update(payment); err != nil {
		return err
	}

	// باز کردن قفل سبد خرید
	return p.unlockCart(payment.UserID)
}
