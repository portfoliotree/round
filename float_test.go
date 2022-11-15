package round_test

import (
	"fmt"
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/portfoliotree/round"
)

func TestRecursive(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		o := NewWithT(t)

		float := 9.87654321

		err := round.Recursive(&float, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(float).To(Equal(9.877))
	})

	t.Run("float64 ptr", func(t *testing.T) {
		o := NewWithT(t)

		float := 9.87654321
		data := &float

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(float).To(Equal(9.877))
	})

	t.Run("field", func(t *testing.T) {
		o := NewWithT(t)

		data := struct {
			Float float64
		}{
			Float: 9.87654321,
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data.Float).To(Equal(9.877))
	})

	t.Run("slice", func(t *testing.T) {
		o := NewWithT(t)

		data := []float64{
			1.87654321,
			2.87654321,
			3.87654321,
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data).To(Equal([]float64{
			1.877,
			2.877,
			3.877,
		}))
	})

	t.Run("array", func(t *testing.T) {
		o := NewWithT(t)

		const size = 3
		data := [size]float64{
			1.87654321,
			2.87654321,
			3.87654321,
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data).To(Equal([size]float64{
			1.877,
			2.877,
			3.877,
		}))
	})

	t.Run("map values", func(t *testing.T) {
		o := NewWithT(t)

		data := map[string]float64{
			"a": 1.87654321,
			"b": 2.87654321,
			"c": 3.87654321,
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data).To(Equal(map[string]float64{
			"a": 1.877,
			"b": 2.877,
			"c": 3.877,
		}))
	})

	t.Run("private field", func(t *testing.T) {
		o := NewWithT(t)

		data := struct {
			float float64
		}{
			float: 9.87654321,
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data.float).To(Equal(9.87654321))
	})

	t.Run("slice field", func(t *testing.T) {
		o := NewWithT(t)

		data := struct {
			Floats []float64
		}{
			Floats: []float64{
				1.87654321,
				2.87654321,
				3.87654321,
			},
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data.Floats).To(Equal([]float64{
			1.877,
			2.877,
			3.877,
		}))
	})

	t.Run("field with precision tag", func(t *testing.T) {
		o := NewWithT(t)

		n := 5.4321
		data := struct {
			Floats []float64 `precision:"1"`
			Ptr    *float64  `precision:"0"`
		}{
			Floats: []float64{
				1.87654321,
				2.87654321,
				3.87654321,
			},
			Ptr: &n,
		}

		err := round.Recursive(&data, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data.Floats).To(Equal([]float64{
			1.9,
			2.9,
			3.9,
		}))
		o.Expect(*data.Ptr).To(Equal(float64(5)))
	})

	t.Run("math.Inf", func(t *testing.T) {
		o := NewWithT(t)

		float := math.Inf(1)

		err := round.Recursive(&float, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(float).To(Equal(math.Inf(1)),
			"it should not be changed")
	})

	t.Run("negative math.Inf", func(t *testing.T) {
		o := NewWithT(t)

		float := math.Inf(-1)

		err := round.Recursive(&float, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(float).To(Equal(math.Inf(-1)),
			"it should not be changed")
	})

	t.Run("negative math.Inf", func(t *testing.T) {
		o := NewWithT(t)

		float := math.Inf(-1)

		err := round.Recursive(&float, 3)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(float).To(Equal(math.Inf(-1)),
			"it should not be changed")
	})

	t.Run("handles errors", func(t *testing.T) {
		o := NewWithT(t)

		f := math.NaN()
		data := struct {
			Floats []*float64
		}{
			Floats: []*float64{
				&f,
			},
		}

		err := round.Recursive(&data, 0)
		o.Expect(err).To(MatchError(
			And(
				ContainSubstring("not a number"),
				ContainSubstring(".Floats[0]"),
			),
		))
	})

	t.Run("malformed precision tag", func(t *testing.T) {
		o := NewWithT(t)

		var data struct {
			Float float64 `precision:"abc"`
		}

		err := round.Recursive(&data, 0)
		o.Expect(err).To(MatchError(
			ContainSubstring("precision tag"),
		))
	})

	t.Run("negative precision tag", func(t *testing.T) {
		o := NewWithT(t)

		var data struct {
			Float float64 `precision:"-1"`
		}
		data.Float = 987.654

		err := round.Recursive(&data, 0)
		o.Expect(err).NotTo(HaveOccurred())

		o.Expect(data.Float).To(Equal(float64(990)))
	})
}

func ExampleDecimal() {
	const pi = 3.14159265359

	// rounding up the 4th decimal place
	piToFourDigits := round.Decimal(pi, 4)
	fmt.Printf("%g == %g is %t\n", piToFourDigits, 3.1416, piToFourDigits == 3.1416)
	fmt.Printf("%g < %g is %t\n", piToFourDigits, pi, piToFourDigits < 3.1416)

	// rounding down the 2nd decimal place
	piToTwoDigits := round.Decimal(pi, 2)
	fmt.Printf("%g < %g is %t\n", piToTwoDigits, pi, piToTwoDigits < pi)

	fmt.Printf("%q\n", fmt.Sprintf("%g", piToTwoDigits))

	// Output: 3.1416 == 3.1416 is true
	// 3.1416 < 3.14159265359 is false
	// 3.14 < 3.14159265359 is true
	// "3.14"
}
