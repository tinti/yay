package main

import (
	"fmt"
	"syscall"
	"unicode"
)

const gitEmptyTree = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

func stringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func removeInvalidTargets(targets []string) []string {
	filteredTargets := make([]string, 0)

	for _, target := range targets {
		db, _ := splitDBFromName(target)

		if db == "aur" && mode == modeRepo {
			fmt.Printf("%s %s %s\n", bold(yellow(arrow)), cyan(target), bold("Can't use target with option --repo -- skipping"))
			continue
		}

		if db != "aur" && db != "" && mode == modeAUR {
			fmt.Printf("%s %s %s\n", bold(yellow(arrow)), cyan(target), bold("Can't use target with option --aur -- skipping"))
			continue
		}

		filteredTargets = append(filteredTargets, target)
	}

	return filteredTargets
}

// LessRunes compares two rune values, and returns true if the first argument is lexicographicaly smaller.
func LessRunes(iRunes, jRunes []rune) bool {
	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}

	return len(iRunes) < len(jRunes)
}

const (
	IOPRIO_CLASS_NONE = iota
	IOPRIO_CLASS_RT
	IOPRIO_CLASS_BE
	IOPRIO_CLASS_IDLE
)

const (
	IOPRIO_WHO_PROCESS = iota + 1
	IOPRIO_WHO_PGRP
	IOPRIO_WHO_USER
)

const IOPRIO_CLASS_SHIFT = 13
const IOPRIO_PRIO_MASK = (1 << IOPRIO_CLASS_SHIFT) - 1

func IOPRIO_PRIO_CLASS(mask int) int        { return (mask) >> IOPRIO_CLASS_SHIFT }
func IOPRIO_PRIO_DATA(mask int) int         { return ((mask) & IOPRIO_PRIO_MASK) }
func IOPRIO_PRIO_VALUE(class, data int) int { return ((class) << IOPRIO_CLASS_SHIFT) | data }

func ioPrioSet(which, who, ioprio int) int {
	ecode, _, _ := syscall.Syscall(syscall.SYS_IOPRIO_SET, uintptr(which), uintptr(who), uintptr(ioprio))
	return int(ecode)
}

const (
	PRIO_PROCESS = iota
	PRIO_PGRP
	PRIO_USER
)

func setPriority(which, who, prio int) int {
	ecode, _, _ := syscall.Syscall(syscall.SYS_SETPRIORITY, uintptr(which), uintptr(who), uintptr(prio))
	return int(ecode)
}
