package tests

import (
	"testing"

	"paypath/internal/liquid"
	"paypath/internal/services/auth"
	"paypath/internal/services/debts"
	"paypath/internal/services/expenses"
	"paypath/internal/services/income"
	"paypath/internal/storage"
)

func TestExpensesRepoCRUD(t *testing.T) {
	repo := expenses.NewRepository(storage.NewTest(t))

	created, err := repo.Create(testUserID, expenses.Expense{Expense: "rent", Cost: 1000, Frequency: "monthly"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("expected an assigned ID")
	}

	list, err := repo.All(testUserID)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("want 1 expense, got %d", len(list))
	}

	upd, err := repo.Update(testUserID, created.ID, expenses.Expense{Cost: 1500})
	if err != nil {
		t.Fatal(err)
	}
	if upd == nil || !approx(upd.Cost, 1500, 0.001) {
		t.Fatalf("update failed: %+v", upd)
	}

	ok, err := repo.Delete(testUserID, created.ID)
	if err != nil || !ok {
		t.Fatalf("delete failed: ok=%v err=%v", ok, err)
	}
}

func TestDebtsRepoCRUD(t *testing.T) {
	repo := debts.NewRepository(storage.NewTest(t))

	created, err := repo.Create(testUserID, debts.Debt{Name: "cc", Type: "credit_card", APY: 20, Balance: 1000})
	if err != nil {
		t.Fatal(err)
	}
	ok, err := repo.Update(testUserID, created.ID, debts.Debt{Name: "cc", Type: "credit_card", APY: 18, Balance: 900})
	if err != nil || !ok {
		t.Fatalf("update failed: ok=%v err=%v", ok, err)
	}
	ok, err = repo.Delete(testUserID, created.ID)
	if err != nil || !ok {
		t.Fatalf("delete failed: ok=%v err=%v", ok, err)
	}
}

func TestIncomeRepoCRUD(t *testing.T) {
	repo := income.NewRepository(storage.NewTest(t))

	created, err := repo.Create(testUserID, income.Income{Job: "Vet", PayType: "salary", AnnualSalary: fptr(60000)})
	if err != nil {
		t.Fatal(err)
	}
	upd, err := repo.Update(testUserID, created.ID, income.Income{Job: "Vet Assistant"})
	if err != nil {
		t.Fatal(err)
	}
	if upd == nil || upd.Job != "Vet Assistant" {
		t.Fatalf("update failed: %+v", upd)
	}
	ok, err := repo.Delete(testUserID, created.ID)
	if err != nil || !ok {
		t.Fatalf("delete failed: ok=%v err=%v", ok, err)
	}
}

func TestLiquidRepoCRUD(t *testing.T) {
	repo := liquid.NewRepository(storage.NewTest(t))

	created, err := repo.Create(testUserID, liquid.Liquid{Bank: "Chase", Balance: 1000})
	if err != nil {
		t.Fatal(err)
	}
	ok, err := repo.Update(testUserID, created.ID, liquid.Liquid{Bank: "Chase", Balance: 2000})
	if err != nil || !ok {
		t.Fatalf("update failed: ok=%v err=%v", ok, err)
	}
	list, err := repo.All(testUserID)
	if err != nil || len(list) != 1 {
		t.Fatalf("want 1 liquid account, got %d (err=%v)", len(list), err)
	}
}

func TestAuthRepoCRUD(t *testing.T) {
	repo := auth.NewRepository(storage.NewTest(t))

	id, err := repo.CreateUser("a@b.com", "hash", "Alice")
	if err != nil {
		t.Fatal(err)
	}
	u, err := repo.GetUserByEmail("a@b.com")
	if err != nil || u == nil || u.ID != int(id) {
		t.Fatalf("GetUserByEmail failed: %+v err=%v", u, err)
	}
	if u2, err := repo.GetUserByID(int(id)); err != nil || u2 == nil {
		t.Fatalf("GetUserByID failed: %+v err=%v", u2, err)
	}

	if rev, _ := repo.IsTokenRevoked("tok123"); rev {
		t.Fatal("token should not be revoked initially")
	}
	if err := repo.RevokeToken("tok123"); err != nil {
		t.Fatal(err)
	}
	if rev, _ := repo.IsTokenRevoked("tok123"); !rev {
		t.Fatal("token should be revoked after RevokeToken")
	}

	if ok, err := repo.DeleteUser(int(id)); err != nil || !ok {
		t.Fatalf("DeleteUser failed: ok=%v err=%v", ok, err)
	}
}
