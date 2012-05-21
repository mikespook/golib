package autoinc

import (
    "testing"
)

func TestAutoInc(t *testing.T) {
    ai := New(0, 2)
    defer ai.Close()
    for i := 0; i < 10; i = i + 2 {
        if id := ai.Id(); id != i {
            t.Errorf("The id 0 wanted, but got %d", id)
        }
    }

}
