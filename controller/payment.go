package controller

import (
	"bm/db"
	"bm/di"
	"bm/entity"
	"bm/repository"
	"errors"
	"fmt"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/precheckoutquery"
)

type PaymentController struct {
	BaseController
	paymentRepo *repository.PaymentRepo
	em          db.Saver
}

func NewPaymentController(
	bc BaseController,
	em db.Saver,
	paymentRepo *repository.PaymentRepo,
) *PaymentController {
	return &PaymentController{
		BaseController: bc,
		paymentRepo:    paymentRepo,
		em:             em,
	}
}

func (c *PaymentController) Premium(b *gotgbot.Bot, ctx *ext.Context) error {
	user := c.User(ctx)
	payment := &entity.Payment{
		UserID:   user.ID,
		Currency: "XTR",
		Price:    100,
		State:    entity.PendingPaymentState,
	}
	c.em.Save(payment)

	_, err := b.SendInvoice(
		ctx.EffectiveSender.Id(),
		"Product name",
		"Product description",
		fmt.Sprintf("%d", payment.ID),
		payment.Currency,
		[]gotgbot.LabeledPrice{{Label: "one", Amount: 100}},
		&gotgbot.SendInvoiceOpts{ProtectContent: true},
	)

	return err
}

func (c *PaymentController) PreCheckoutQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	id, err := strconv.Atoi(ctx.PreCheckoutQuery.InvoicePayload)
	if err != nil {
		ctx.PreCheckoutQuery.Answer(b, false, nil)
		return err
	}

	payment := c.paymentRepo.FindByID(uint(id))
	if payment == nil {
		ctx.PreCheckoutQuery.Answer(b, false, nil)

		return errors.New("payment not exists")
	}

	if payment.State != entity.PendingPaymentState {
		ctx.PreCheckoutQuery.Answer(b, false, nil)

		return errors.New("payment not pending")
	}

	// TODO: validate price end user if required

	ctx.PreCheckoutQuery.Answer(b, true, nil)

	return nil
}

func (c *PaymentController) SuccessfulPayment(b *gotgbot.Bot, ctx *ext.Context) error {
	id, err := strconv.Atoi(ctx.EffectiveMessage.SuccessfulPayment.InvoicePayload)
	if err != nil {
		return err
	}

	payment := c.paymentRepo.FindByID(uint(id))
	if payment == nil {
		return errors.New("payment not exists")
	}

	payment.State = entity.SuccessPaymentState
	payment.Currency = ctx.EffectiveMessage.SuccessfulPayment.Currency
	payment.TotalAmount = ctx.EffectiveMessage.SuccessfulPayment.TotalAmount
	payment.InvoicePayload = ctx.EffectiveMessage.SuccessfulPayment.InvoicePayload
	payment.ShippingOptionID = ctx.EffectiveMessage.SuccessfulPayment.ShippingOptionId
	payment.TelegramPaymentChargeID = ctx.EffectiveMessage.SuccessfulPayment.TelegramPaymentChargeId
	payment.ProviderPaymentChargeID = ctx.EffectiveMessage.SuccessfulPayment.ProviderPaymentChargeId

	c.em.Save(payment)

	// TODO: activate payment here

	_, err = ctx.EffectiveMessage.Reply(
		b,
		c.Trans(ctx, "payment.complete.success"),
		nil,
	)

	return err
}

func (c *PaymentController) Register(bot *gotgbot.Bot, dispatcher *ext.Dispatcher) error {
	// NOTE: try /premium command for test
	dispatcher.AddHandler(handlers.NewCommand("premium", c.Premium))

	dispatcher.AddHandler(handlers.NewPreCheckoutQuery(
		precheckoutquery.All,
		c.PreCheckoutQuery,
	))

	dispatcher.AddHandler(handlers.NewMessage(
		message.SuccessfulPayment,
		c.SuccessfulPayment,
	))

	return nil
}

var _ di.Controller = &PaymentController{}
