package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_assert_eq_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	type T struct {
		a int
		b string
	}

	p1 := &T{
		a: 100,
		b: "100",
	}

	p2 := p1

	a1 := [3]int{100, 200, 300}
	a2 := [3]int{100, 200, 300}

	// When / Then
	assert.That(1).IsEq(1)
	assert.That("1").IsEq("1")
	assert.That(1.0).IsEq(1.0)
	assert.That(true).IsEq(true)
	assert.That(false).IsEq(false)

	assert.That(p1).IsEq(p2)

	assert.That(*p1).IsEq(T{a: 100, b: "100"})

	assert.That(T{a: 100, b: "100"}).IsEq(T{a: 100, b: "100"})

	assert.That(a1).IsEq(a2)
}

func Test_assert_deep_eq_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	type SubT struct {
		d map[string]int
	}

	type T struct {
		SubT
		a int
		b string
		c []string
	}

	p1 := &T{
		SubT: SubT{
			d: map[string]int{"100": 100, "200": 200, "300": 300},
		},
		a: 100,
		b: "100",
		c: []string{"100", "200", "300"},
	}

	p2 := p1

	s1 := []int{100, 200, 300}
	s2 := []int{100, 200, 300}

	m1 := map[string]int{"100": 100, "200": 200, "300": 300}
	m2 := map[string]int{"100": 100, "200": 200, "300": 300}

	// When / Then
	assert.That(1).IsDeepEq(1)
	assert.That("1").IsDeepEq("1")
	assert.That(1.0).IsDeepEq(1.0)
	assert.That(true).IsDeepEq(true)
	assert.That(false).IsDeepEq(false)

	assert.That(p1).IsDeepEq(p2)

	assert.That(*p1).IsDeepEq(T{SubT: SubT{d: map[string]int{"100": 100, "200": 200, "300": 300}},
		a: 100, b: "100", c: []string{"100", "200", "300"}})
	assert.That(T{SubT: SubT{d: map[string]int{"100": 100, "200": 200, "300": 300}},
		a: 100, b: "100", c: []string{"100", "200", "300"}}).
		IsDeepEq(T{SubT: SubT{d: map[string]int{"100": 100, "200": 200, "300": 300}},
			a: 100, b: "100", c: []string{"100", "200", "300"}})

	assert.That(T{a: 100, b: "100"}).IsDeepEq(T{a: 100, b: "100"})

	assert.That(s1).IsDeepEq(s2)
	assert.That(m1).IsDeepEq(m2)
}

func Test_assert_eq_negation_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	type T struct {
		a int
		b string
	}

	p1 := &T{
		a: 100,
		b: "100",
	}

	p2 := &T{
		a: 100,
		b: "100",
	}

	a1 := [3]int{100, 200, 300}
	a2 := [3]int{100, 200, 400}

	// When / Then
	assert.That(1).Not().IsEq(2)
	assert.That("1").Not().IsEq("2")
	assert.That(1.0).Not().IsEq(1.2)
	assert.That(true).Not().IsEq(false)
	assert.That(false).Not().IsEq(true)

	assert.That(p1).Not().IsEq(p2)

	assert.That(*p1).Not().IsEq(T{a: 100, b: "101"})

	assert.That(T{a: 100, b: "100"}).Not().IsEq(T{a: 100, b: "101"})

	assert.That(a1).Not().IsEq(a2)
}

func Test_assert_deep_eq_negation_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	type SubT struct {
		d map[string]int
	}

	type T struct {
		SubT
		a int
		b string
		c []string
	}

	p1 := &T{
		SubT: SubT{
			d: map[string]int{"100": 100, "200": 200, "300": 300},
		},
		a: 100,
		b: "100",
		c: []string{"100", "200", "300"},
	}

	p2 := &T{
		SubT: SubT{
			d: map[string]int{"100": 100, "200": 200, "300": 300},
		},
		a: 100,
		b: "100",
		c: []string{"100", "200", "400"},
	}

	s1 := []int{100, 200, 300}
	s2 := []int{100, 200, 400}

	m1 := map[string]int{"100": 100, "200": 200, "300": 300}
	m2 := map[string]int{"100": 100, "200": 200, "300": 400}

	// When / Then
	assert.That(1).Not().IsDeepEq(2)
	assert.That("1").Not().IsDeepEq("2")
	assert.That(1.0).Not().IsDeepEq(1.2)
	assert.That(true).Not().IsDeepEq(false)
	assert.That(false).Not().IsDeepEq(true)

	assert.That(p1).Not().IsDeepEq(p2)

	assert.That(*p1).Not().IsDeepEq(T{SubT: SubT{d: map[string]int{"100": 100, "200": 200, "300": 300}},
		a: 100, b: "100", c: []string{"100", "200", "400"}})
	assert.That(T{SubT: SubT{d: map[string]int{"100": 100, "200": 200, "300": 300}},
		a: 100, b: "100", c: []string{"100", "200", "300"}}).Not().
		IsDeepEq(T{SubT: SubT{d: map[string]int{"100": 100, "200": 200, "300": 300}},
			a: 100, b: "100", c: []string{"100", "200", "400"}})

	assert.That(T{a: 100, b: "100"}).Not().IsDeepEq(T{a: 100, b: "200"})

	assert.That(s1).Not().IsDeepEq(s2)
	assert.That(m1).Not().IsDeepEq(m2)
}

func Test_assert_eq_should_not_pass(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	testEntries := []struct {
		a           interface{}
		b           interface{}
		expectation func(*mocks.MockPublicTB)
	}{
		{
			a: 1,
			b: 2,
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.EXPECT().Errorf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v", 2, 1)
			},
		},
		{
			a: "1",
			b: "2",
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.EXPECT().Errorf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v", "2", "1")
			},
		},
		{
			a: [3]int{100, 200, 300},
			b: [3]int{100, 200, 400},
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.
					EXPECT().
					Errorf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v",
						[3]int{100, 200, 400}, [3]int{100, 200, 300})
			},
		},
	}

	for _, entry := range testEntries {
		// Given
		tMock := mocks.NewMockPublicTB(ctrl)
		assert := assertion.New(tMock)

		// Expectation
		tMock.EXPECT().Helper().AnyTimes()
		entry.expectation(tMock)

		// When
		assert.That(entry.a).IsEq(entry.b)
	}
}

func Test_assert_deep_eq_should_not_pass(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	testEntries := []struct {
		a           interface{}
		b           interface{}
		expectation func(*mocks.MockPublicTB)
	}{
		{
			a: 1,
			b: 2,
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.EXPECT().Errorf("\nValue is not deep equal to expectation.\nExpected : %v\nGot : %v", 2, 1)
			},
		},
		{
			a: "1",
			b: "2",
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.EXPECT().Errorf("\nValue is not deep equal to expectation.\nExpected : %v\nGot : %v", "2", "1")
			},
		},
		{
			a: [3]int{100, 200, 300},
			b: [3]int{100, 200, 400},
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.
					EXPECT().
					Errorf("\nValue is not deep equal to expectation.\nExpected : %v\nGot : %v",
						[3]int{100, 200, 400}, [3]int{100, 200, 300})
			},
		},
	}

	for _, entry := range testEntries {
		// Given
		tMock := mocks.NewMockPublicTB(ctrl)
		assert := assertion.New(tMock)

		// Expectation
		tMock.EXPECT().Helper().AnyTimes()
		entry.expectation(tMock)

		// When
		assert.That(entry.a).IsDeepEq(entry.b)
	}
}
func Test_assert_eq_negation_should_not_pass(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	testEntries := []struct {
		a           interface{}
		b           interface{}
		expectation func(*mocks.MockPublicTB)
	}{
		{
			a: 1,
			b: 1,
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.EXPECT().Errorf("\nValue should not be equal to : %v", 1)
			},
		},
		{
			a: "1",
			b: "1",
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.EXPECT().Errorf("\nValue should not be equal to : %v", "1")
			},
		},
		{
			a: [3]int{100, 200, 300},
			b: [3]int{100, 200, 300},
			expectation: func(tMock *mocks.MockPublicTB) {
				tMock.
					EXPECT().
					Errorf("\nValue should not be equal to : %v", [3]int{100, 200, 300})
			},
		},
	}

	for _, entry := range testEntries {
		// Given
		tMock := mocks.NewMockPublicTB(ctrl)
		assert := assertion.New(tMock)

		// Expectation
		tMock.EXPECT().Helper().AnyTimes()
		entry.expectation(tMock)

		// When
		assert.That(entry.a).Not().IsEq(entry.b)
	}
}
func Test_assert_eq_should_not_pass_with_custom_message(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	tMock := mocks.NewMockPublicTB(ctrl)
	assert := assertion.New(tMock)

	// Expectation
	tMock.EXPECT().Helper().AnyTimes()
	tMock.EXPECT().Errorf("This is custom message")

	// When
	assert.That(1).Logf("This is custom message").IsEq(2)
}

func Test_assert_eq_should_not_pass_and_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	tMock := mocks.NewMockPublicTB(ctrl)
	assert := assertion.New(tMock)

	// Expectation
	tMock.EXPECT().Helper().AnyTimes()
	tMock.EXPECT().Fatalf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v", 2, 1)

	// When
	assert.That(1).OrFatal().IsEq(2)
}

func Test_assert_eq_should_not_pass_silently(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	tMock := mocks.NewMockPublicTB(ctrl)
	assert := assertion.New(tMock)

	// Expectation
	tMock.EXPECT().Helper().AnyTimes()
	tMock.EXPECT().Fail()

	// When
	assert.That(1).Silent().IsEq(2)
}

func Test_assert_eq_should_not_pass_silently_and_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	tMock := mocks.NewMockPublicTB(ctrl)
	assert := assertion.New(tMock)

	// Expectation
	tMock.EXPECT().Helper().AnyTimes()
	tMock.EXPECT().FailNow()

	// When
	assert.That(1).OrFatal().Silent().IsEq(2)
}
