package main

import (
	"context"
	"os"

	"github.com/raishey/plinko"
	"github.com/raishey/plinko/pkg/config"
	"github.com/raishey/plinko/pkg/config/state"
	"github.com/raishey/plinko/pkg/renderers"
)

const Created plinko.State = "Created"
const Opened plinko.State = "Opened"
const Claimed plinko.State = "Claimed"
const ArriveAtStore plinko.State = "ArrivedAtStore"
const MarkedAsPickedUp plinko.State = "MarkedAsPickedup"
const Delivered plinko.State = "Delivered"
const Canceled plinko.State = "Canceled"
const Returned plinko.State = "Returned"

const Submit plinko.Trigger = "Submit"
const AddItemToOrder plinko.Trigger = "AddItemToOrder"
const Cancel plinko.Trigger = "Cancel"
const Open plinko.Trigger = "Open"
const Claim plinko.Trigger = "Claim"
const Deliver plinko.Trigger = "Deliver"
const Return plinko.Trigger = "Return"
const Reinstate plinko.Trigger = "Reinstate"

func OnNewOrderEntry(ctx context.Context, p plinko.Payload, t plinko.TransitionInfo) (plinko.Payload, error) {
	// perform a series of steps based on the
	// payload and transition info
	// ...

	return p, nil
}

func main() {
	p := config.CreatePlinkoDefinition()

	p.Configure(Created, state.WithDescription("Initially submitted")).
		OnEntry(OnNewOrderEntry).
		Permit(Open, Opened).
		Permit(Cancel, Canceled)

	p.Configure(Opened).
		Permit(AddItemToOrder, Opened).
		Permit(Claim, Claimed).
		Permit(Cancel, Canceled)

	p.Configure(Claimed).
		Permit(AddItemToOrder, Claimed).
		Permit(Submit, ArriveAtStore).
		Permit(Cancel, Canceled)

	p.Configure(ArriveAtStore).
		Permit(Submit, MarkedAsPickedUp).
		Permit(Cancel, Canceled)

	p.Configure(MarkedAsPickedUp).
		Permit(Deliver, Delivered).
		Permit(Cancel, Canceled)

	p.Configure(Delivered).
		Permit(Return, Returned)

	p.Configure(Canceled).
		Permit(Reinstate, Created)

	p.Configure(Returned, state.WithDescription("An order is Returned if it has been returned to the store.")).Permit(Cancel, Canceled)

	compilerOutput := p.Compile()

	for _, m := range compilerOutput.Messages {
		if m.CompileMessage == plinko.CompileError {
			panic("FSM compilation error: " + m.CompileMessage)
		}
	}

	umlFile, err := os.OpenFile("plinko.uml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic("failed to open file plinko.uml: " + err.Error())
	}
	defer umlFile.Close()

	uml := renderers.NewUML(umlFile)

	dotFile, err := os.OpenFile("plinko.dot", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic("failed to open file plinko.dot" + err.Error())
	}
	defer dotFile.Close()

	dot := renderers.NewDot(dotFile)

	p.Render(dot)
	p.Render(uml)
}
