package main

import "fmt"

//go:generate stringer -type=EnrollmentStatus --trimprefix EnrollmentStatus

type EnrollmentStatus int

const (
	EnrollmentStatusPotential EnrollmentStatus = iota
	EnrollmentStatusEnrolled
	EnrollmentStatusGraduated

	EnrollmentStatusWithdraw
	EnrollmentStatusNonPotential
	EnrollmentStatusLOA
)

/*func (enrollmentStatus EnrollmentStatus) String() string {
	return "EnrollmentStatus"
}*/

func main() {
	fmt.Println()
	fmt.Println(EnrollmentStatusLOA)
}
