package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/flow"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Evaluator interface {
	String() string
	EvalPar(*Span, *Func, Arg) (Return, Effect, error)
	EvalSeq(*Span, *Func, Arg) (Return, Effect, error)
}

// Program is an Evaluator.
type Program struct {
	Idiom  Repo `ko:"name=idiom"`
	Repo   Repo `ko:"name=repo"` // function repository from user sources
	System `ko:"name=system"`
}

type System struct {
	Faculty  Faculty  `ko:"name=faculty"`  // operator logic
	Boundary Boundary `ko:"name=boundary"` // circuit surgery
	Combiner Combiner `ko:"name=combiner"` // circuit surgery
}

func (prog Program) String() string { return Sprint(prog) }

func (f evalFlow) StepResidue() *StepResidue {
	return &StepResidue{Frame: f.Frame, Shape: f.Shape, Effect: f.Effect}
}

func FlowResidues(stepFlow []Flow) (stepResidue []*StepResidue) {
	stepResidue = make([]*StepResidue, len(stepFlow))
	for i, f := range stepFlow {
		stepResidue[i] = f.(evalFlow).StepResidue()
	}
	return
}

// func (prog Program) Eval(span *Span, f *Func, arg Arg) (Return, Effect, error) {
// 	return prog.EvalSeq(span, f, arg)
// }

func (prog Program) EvalSeq(span *Span, f *Func, arg Arg) (Return, Effect, error) {
	envelope := evalEnvelope{Program: prog, Arg: arg, Span: span}
	if returnFlow, stepFlow, err := PlaySeqFlow(span, f, envelope); err != nil {
		return nil, nil, err
	} else {
		returned := returnFlow.(evalFlow).Return
		if eff, err := prog.Combiner.Combine(span, f, arg, returned, FlowResidues(stepFlow)); err != nil {
			return nil, nil, err
		} else {
			return returned, eff, nil
		}
	}
}

func (prog Program) EvalPar(span *Span, f *Func, arg Arg) (Return, Effect, error) {
	envelope := evalEnvelope{Program: prog, Arg: arg, Span: span}
	if returnFlow, stepFlow, err := PlayParFlow(span, f, envelope); err != nil {
		return nil, nil, err
	} else {
		returned := returnFlow.(evalFlow).Return
		if eff, err := prog.Combiner.Combine(span, f, arg, returned, FlowResidues(stepFlow)); err != nil {
			return nil, nil, err
		} else {
			return returned, eff, nil
		}
	}
}