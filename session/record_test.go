package session

import "testing"

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	table := NewSession().Model(&User{})
	err1 := table.DropTable()
	err2 := table.CreateTable()
	_, err3 := table.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return table
}

func TestSession_Insert(t *testing.T) {
	table := testRecordInit(t)
	affected, err := table.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	table := testRecordInit(t)
	var rows []User
	err := table.Find(&rows)
	if err != nil || len(rows) != 2 {
		t.Fatalf("failed to query all, err:%v, len:%d", err, len(rows))
	}
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count() // 同一个 session，执行之后会清空 sql，但是 clause 并不会清空

	if affected != 1 || count != 0 {
		t.Fatal("failed to delete or count")
	}
}
