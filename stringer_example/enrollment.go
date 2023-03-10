package main

import "fmt"

//go:generate stringer -type=EnrollmentStatus --trimprefix EnrollmentStatus

type EnrollmentStatus int

const (
	EnrollmentStatusPotential EnrollmentStatus = iota
	EnrollmentStatusEnrolled

	EnrollmentStatusNonPotential
	EnrollmentStatusLOA
	EnrollmentStatusGraduated
	EnrollmentStatusWithdraw
)

/*func (enrollmentStatus EnrollmentStatus) String() string {
	//...
	//...
	return "EnrollmentStatus"
}*/

func main() {
	fmt.Println()
	fmt.Println(EnrollmentStatusPotential)
}
