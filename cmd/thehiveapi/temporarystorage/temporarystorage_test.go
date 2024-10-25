package temporarystoarge

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemporaryStorage(t *testing.T) {
	var (
		ts *TemporaryStorage

		err error

		test_uniq_case_id_1 []byte = []byte("uniq_case_id:f78773r88r8w874et7rt7g77sw7w")
		test_uniq_case_id_2 []byte = []byte("uniq_case_id:g7627gdff8fyr8298euihusd8y823")
		test_uniq_case_id_3 []byte = []byte("uniq_case_id:fs662te73t73tr73t6rt37tr7376r3")

		test_uuid_1 string
		test_uuid_2 string
		test_uuid_3 string
	)

	ts, err = NewTemporaryStorage(10)
	assert.NoError(t, err)

	test_uuid_1 = ts.SetValue("case_id_1", test_uniq_case_id_1)
	d, ok := ts.GetValue(test_uuid_1)
	assert.True(t, ok)
	assert.Equal(t, test_uniq_case_id_1, d)

	test_uuid_2 = ts.SetValue("case_id_2", test_uniq_case_id_2)
	d, ok = ts.GetValue(test_uuid_2)
	assert.True(t, ok)
	assert.Equal(t, test_uniq_case_id_2, d)

	//удаление
	ts.DeleteElement(test_uuid_1)
	_, ok = ts.GetValue(test_uuid_1)
	assert.False(t, ok)

	//ставим паузу для автоматического удаления устаревших значений
	time.Sleep(9 * time.Second)

	test_uuid_3 = ts.SetValue("case_id_3", test_uniq_case_id_3)

	time.Sleep(6 * time.Second)
	//удаляется автоматически
	_, ok = ts.GetValue(test_uuid_2)
	assert.False(t, ok)

	d, ok = ts.GetValue(test_uuid_3)
	assert.True(t, ok)
	assert.Equal(t, test_uniq_case_id_3, d)
}
