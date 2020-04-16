package queryutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleQuery(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "code",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " code = ? "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestTwoQuery(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "company__code",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " company_id in ( select id from trinity_company where code = ? ) "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestThreeQuery(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "costcenter__company__code",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code = ? ) ) "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourQuery(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__code",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code = ? ) ) ) "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestIlike(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "ilike",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " ilike = ? "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestSingleIlike(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "code__ilike",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " code ilike ? "
	expectValueSQL := "%test%"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}
func TestFourIlike(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__code__ilike",
		QueryValue:  "test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code ilike ? ) ) ) "
	expectValueSQL := "%test%"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourIn(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__code__in",
		QueryValue:  "test,test,test",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code in (?) ) ) ) "
	expectValueSQL := []string{"test", "test", "test"}

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourStart(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__create_time__start",
		QueryValue:  "2019-01-01",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time >= ? ) ) ) "
	expectValueSQL := "2019-01-01 00:00:00"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourEnd(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__create_time__end",
		QueryValue:  "2019-01-01",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time <= ? ) ) ) "
	expectValueSQL := "2019-01-01 23:59:59"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourIsNullFalse(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__create_time__isnull",
		QueryValue:  "false",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time is not null ) ) ) "
	// expectValueSQL := nil

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, nil, a.ValueSQL, "value error")

}

func TestFourIsNullTrue(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__create_time__isnull",
		QueryValue:  "true",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time is null ) ) ) "
	// expectValueSQL := nil

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, nil, a.ValueSQL, "value error")

}

func TestFourIsEmptyTrue(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__create_time__isempty",
		QueryValue:  "true",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where (COALESCE(\"create_time\"::varchar ,'') = '' ) ) ) ) "
	// expectValueSQL := nil

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, nil, a.ValueSQL, "value error")

}

func TestFourIsEmptyFalse(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__create_time__isempty",
		QueryValue:  "false",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where (COALESCE(\"create_time\"::varchar ,'') != '' ) ) ) ) "
	// expectValueSQL := nil

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, nil, a.ValueSQL, "value error")

}

func TestFourLT(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__qty__lt",
		QueryValue:  "50",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty < ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourLTE(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__qty__lte",
		QueryValue:  "50",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty <= ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourGT(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__qty__gt",
		QueryValue:  "50",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty > ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}

func TestFourGTE(t *testing.T) {
	a := &FilterQuery{
		QueryName:   "asset__costcenter__company__qty__gte",
		QueryValue:  "50",
		TablePrefix: "trinity_",
	}
	a.GetFilterQuerySQL()
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty >= ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, a.ConditionSQL, "conditionsql error")
	assert.Equal(t, expectValueSQL, a.ValueSQL, "value error")
}
