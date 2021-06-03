package processes

import (
	"amaterasu/cqrs/infrastructure/messaging"
	"amaterasu/utils"
	"log"
	"reflect"
)

const (
	NotStarted ProcessState = iota
	AwaitingReservationConfirmation
	ReservationConfirmationReceived
	PaymentConfirmationReceived
)

type (
	ProcessState uint

	OrderPlaced struct {
		messaging.Event
	}

	OrderUpdated struct {
		messaging.Event
	}

	OrderConfirmed struct {
		messaging.Event
	}

	RegistrationProcessmanager struct {
		ProcessManager
		State ProcessState
	}

	RegistrationProcessmanagerRouter struct {
		contextFactory IProcessManagerDataContext
	}
)

func (pm *RegistrationProcessmanager) HandleOrderPlaced(event OrderPlaced) error {

	pm.Completed = true
	log.Println(utils.GetObjType2(reflect.TypeOf(event)))
	return nil
}

func (pm *RegistrationProcessmanager) HandleOrderUpdated(event OrderUpdated) error {
	log.Println(utils.GetObjType2(reflect.TypeOf(event)))
	return nil
}

func (pm *RegistrationProcessmanager) HandleOrderConfirmed(event OrderConfirmed) error {
	pm.Completed = true
	log.Println(utils.GetObjType2(reflect.TypeOf(event)))
	return nil
}

// ============================
func (router RegistrationProcessmanagerRouter) New_RegistrationProcessmanagerRouter(processmanager IProcessManager, contextFactory IProcessManagerDataContext) *RegistrationProcessmanagerRouter {
	return &RegistrationProcessmanagerRouter{
		contextFactory: contextFactory,
	}
}

func (router RegistrationProcessmanagerRouter) castPm(ipm IProcessManager) *RegistrationProcessmanager {
	if ipm == nil {
		pm := &RegistrationProcessmanager{}
		pm.Id = utils.NewUuidString()
		return pm
	}
	return ipm.(*RegistrationProcessmanager)

}

func (router RegistrationProcessmanagerRouter) HandleOrderPlaced(event OrderPlaced) error {

	pm := router.castPm(router.contextFactory.Find(event.SourceId, true))

	if err := pm.HandleOrderPlaced(event); err != nil {
		return err
	}

	if err := router.contextFactory.Save(pm); err != nil {
		return err
	}

	return nil
}

func (router RegistrationProcessmanagerRouter) HandleOrderUpdated(event OrderUpdated) error {

	ipm := router.contextFactory.Find(event.SourceId, true)

	if ipm != nil {

		pm := ipm.(*RegistrationProcessmanager)

		if err := pm.HandleOrderUpdated(event); err != nil {
			return err
		}

		if err := router.contextFactory.Save(pm); err != nil {
			return err
		}

	} else {
		log.Println("Failed to locate the registration process manager handling the order with id ", event.GetSourceId())
	}

	return nil
}

func (router RegistrationProcessmanagerRouter) HandleOrderConfirmed(event OrderConfirmed) error {
	ipm := router.contextFactory.Find(event.SourceId, true)

	if ipm != nil {

		pm := ipm.(*RegistrationProcessmanager)

		if err := pm.HandleOrderConfirmed(event); err != nil {
			return err
		}

		if err := router.contextFactory.Save(pm); err != nil {
			return err
		}

	} else {
		log.Println("Failed to locate the registration process manager to complete with id  ", event.GetSourceId())
	}

	return nil
}
