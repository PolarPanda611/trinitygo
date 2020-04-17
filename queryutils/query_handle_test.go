package queryutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleQuery(t *testing.T) {

	d := NewDecoder("Code", "test", "trinity_")
	expectConditionSQL := " code = ? "
	expectValueSQL := "test"
	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestTwoQuery(t *testing.T) {
	d := NewDecoder("company__code", "test", "trinity_")
	expectConditionSQL := " company_id in ( select id from trinity_company where code = ? ) "
	expectValueSQL := "test"
	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestThreeQuery(t *testing.T) {
	d := NewDecoder("costcenter__company__code", "test", "trinity_")
	expectConditionSQL := " costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code = ? ) ) "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourQuery(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__code", "test", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code = ? ) ) ) "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestIlike(t *testing.T) {
	d := NewDecoder("ilike", "test", "trinity_")
	expectConditionSQL := " ilike = ? "
	expectValueSQL := "test"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestSingleIlike(t *testing.T) {
	d := NewDecoder("Code__ilike", "test", "trinity_")
	expectConditionSQL := " code ilike ? "
	expectValueSQL := "%test%"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}
func TestFourIlike(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__code__ilike", "test", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code ilike ? ) ) ) "
	expectValueSQL := "%test%"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourIn(t *testing.T) {
	d := NewDecoder("Asset__costcenter__company__Code__in", "test,test,test", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where code in (?) ) ) ) "
	expectValueSQL := []string{"test", "test", "test"}

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourStart(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__create_time__start", "2019-01-01", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time >= ? ) ) ) "
	expectValueSQL := "2019-01-01 00:00:00"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourEnd(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__create_time__end", "2019-01-01", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time <= ? ) ) ) "
	expectValueSQL := "2019-01-01 23:59:59"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourIsNullFalse(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__create_time__isnull", "false", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time is not null ) ) ) "
	var expectValueSQL interface{}

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")

}

func TestFourIsNullTrue(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__create_time__isnull", "true", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where create_time is null ) ) ) "
	var expectValueSQL interface{}

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")

}

func TestFourIsEmptyTrue(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__create_time__isempty", "true", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where (COALESCE(\"create_time\"::varchar ,'') = '' ) ) ) ) "
	var expectValueSQL interface{}

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")

}

func TestFourIsEmptyFalse(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__create_time__isempty", "false", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where (COALESCE(\"create_time\"::varchar ,'') != '' ) ) ) ) "
	var expectValueSQL interface{}

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")

}

func TestFourLT(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__qty__lt", "50", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty < ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourLTE(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__qty__lte", "50", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty <= ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}

func TestFourGT(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__qty__gt", "50", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty > ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}
func TestFourGTE(t *testing.T) {
	d := NewDecoder("asset__costcenter__company__qty__gte", "50", "trinity_")
	expectConditionSQL := " asset_id in ( select id from trinity_asset where costcenter_id in ( select id from trinity_costcenter where company_id in ( select id from trinity_company where qty >= ? ) ) ) "
	expectValueSQL := "50"

	assert.Equal(t, expectConditionSQL, d.ConditionSQL(), "conditionsql error")
	assert.Equal(t, expectValueSQL, d.ValueSQL(), "value error")
}
