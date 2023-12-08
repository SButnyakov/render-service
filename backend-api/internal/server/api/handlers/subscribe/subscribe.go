package subscribe

import (
	"backend-api/internal/config"
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type SubscriptionProvider interface {
	GetExpireDate(int64) (*time.Time, error)
	Create(storage.Subscription) error
	Update(subscription storage.Subscription) error
}

type SubscriptionTypeProvider interface {
	GetID(string) (int64, error)
}

type PaymentCreator interface {
	Create(storage.Payment) error
}

type PaymentTypeProvider interface {
	GetID(string) (int64, error)
}

type Response struct {
	resp.Response
}

func New(log *slog.Logger, cfg *config.Config, pCreator PaymentCreator, ptProvider PaymentTypeProvider,
	stProvider SubscriptionTypeProvider, sProvider SubscriptionProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handlers.subscribe.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		uid := r.Context().Value("uid").(int64)

		pTypeId, err := ptProvider.GetID(cfg.Payments.SubPremiumMonth)
		if err != nil {
			log.Error("invalid payment type")
			responseError(w, r, resp.Error("invalid payment type"), http.StatusBadRequest)
			return
		}

		err = pCreator.Create(storage.Payment{
			UserID:   uid,
			TypeId:   pTypeId,
			DateTime: time.Now()})
		if err != nil {
			log.Error("failed to create payment record")
			responseError(w, r, resp.Error("failed to create payment record"), http.StatusInternalServerError)
			return
		}

		sTypeId, err := stProvider.GetID(cfg.Subscriptions.Premium)
		if err != nil {
			log.Error("invalid subscription type")
			responseError(w, r, resp.Error("invalid subscription type"), http.StatusBadRequest)
			return
		}

		expireDate, err := sProvider.GetExpireDate(uid)
		if err != nil {
			if errors.Is(err, storage.ErrSubscriptionNotFound) {
				err = sProvider.Create(storage.Subscription{
					UserId:     uid,
					TypeId:     sTypeId,
					ExpireDate: time.Now().AddDate(0, 1, 0)})
				if err != nil {
					log.Error("failed to subscribe")
					responseError(w, r, resp.Error("subscribing failed"), http.StatusInternalServerError)
					return
				}
				responseOK(w, r)
				return
			} else {
				log.Error("failed to get subscription info")
				responseError(w, r, resp.Error("subscribing failed"), http.StatusInternalServerError)
				return
			}
		}

		err = sProvider.Update(storage.Subscription{
			UserId:     uid,
			TypeId:     sTypeId,
			ExpireDate: expireDate.AddDate(0, 1, 0)})
		if err != nil {
			log.Error("failed to subscribe")
			responseError(w, r, resp.Error("subscribing failed"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
