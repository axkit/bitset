package bitset

import (
	"testing"
)

func TestNew(t *testing.T) {
	bs := New(10)
	if n := bs.Len(); n != 16 {
		t.Errorf("expected bitset length >= 16, got %d", n)
	}
}

func TestByteBitSet_SetAndIsSet(t *testing.T) {
	bs := New(16)
	bs.Set(true, 0, 5, 15)
	if !bs.IsSet(0) || !bs.IsSet(5) || !bs.IsSet(15) {
		t.Error("expected bits 0, 5, and 15 to be set")
	}

	if bs.IsSet(1) {
		t.Error("expected bit 1 to be unset")
	}
	bs.Set(false, 5)
	if bs.IsSet(5) {
		t.Error("expected bit 5 to be cleared")
	}

	if bs.IsSet(1000) {
		t.Error("expected result to be false for out of bounds")
	}
}

func TestByteBitSet_AreSet(t *testing.T) {

	bs, err := ParseHexString("b3")
	if err != nil {
		t.Fatalf("failed to parse string: %v", err)
	}

	t.Run("empty input", func(t *testing.T) {
		res := bs.AreSet(All)
		if res {
			t.Error("expected no bits to be checked")
		}
	})

	t.Run("empty mask", func(t *testing.T) {
		bs := New(0)
		res := bs.AreSet(All, 5)
		if res {
			t.Error("expected no bits to be checked")
		}
	})

	t.Run("valid input: all", func(t *testing.T) {
		res := bs.AreSet(All, 0, 2, 3, 6, 7)
		if !res {
			t.Error("expected bits 0, 2, 3, 6, and 7 to be set")
		}
	})

	t.Run("valid input: all subset", func(t *testing.T) {
		res := bs.AreSet(All, 0, 2)
		if !res {
			t.Error("expected bits 0 and 6 to be set")
		}
	})

	t.Run("valid input: any", func(t *testing.T) {
		res := bs.AreSet(Any, 1, 4, 6)
		if !res {
			t.Error("expected bits 0, 1, or 2 to be set")
		}
	})

	t.Run("valid input: all, invalid subset", func(t *testing.T) {
		res := bs.AreSet(All, 0, 5)
		if res {
			t.Error("expected bit 5 to be unset")
		}
	})

	t.Run("valid input: any, invalid subset", func(t *testing.T) {
		res := bs.AreSet(Any, 5, 15)
		if res {
			t.Error("expected bits 5 and 15 to be unset")
		}
	})

}

func TestByteBitSet_String(t *testing.T) {

	t.Run("empty bit set", func(t *testing.T) {
		bs := New(0)
		s := bs.String()
		if s != "" {
			t.Errorf("expected empty string, got %s", s)
		}
	})
	t.Run("valid bit set", func(t *testing.T) {
		bs := New(8)
		bs.Set(true, 0, 6)
		s := bs.String()
		parsed, err := ParseHexString(s)
		if err != nil {
			t.Fatalf("failed to parse string: %v", err)
		}
		if !parsed.IsSet(0) || !parsed.IsSet(6) {
			t.Error("parsed bitset does not match original")
		}
	})
}

func TestParseBinaryString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		bs, err := ParseBinaryString("")
		if err != nil {
			t.Errorf("unexpected error for empty string: %v", err)
		}
		if bs.Len() != 0 {
			t.Error("expected zero-length bitset")
		}
	})
	t.Run("invalid binary string", func(t *testing.T) {

		_, err := ParseBinaryString("101a")
		if err == nil {
			t.Error("expected error for invalid binary string")
		}
	})
	t.Run("valid binary string", func(t *testing.T) {
		bs, err := ParseBinaryString("10101010")
		if err != nil {
			t.Fatalf("unexpected error for valid binary string: %v", err)
		}

		if bs.Len() != 8 {
			t.Errorf("expected bitset length 8, got %d", bs.Len())
		}
		if !bs.IsSet(0) || !bs.IsSet(2) || !bs.IsSet(4) {
			t.Error("expected bits 0, 2, and 4 to be set")
		}
		if bs.IsSet(1) || bs.IsSet(3) || bs.IsSet(5) {
			t.Error("expected bits 1, 3, and 5 to be unset")
		}
	})
	t.Run("valid binary string with leading zeros", func(t *testing.T) {
		bs, err := ParseBinaryString("00000001")
		if err != nil {
			t.Fatalf("unexpected error for valid binary string: %v", err)
		}
		if bs.Len() != 8 {
			t.Errorf("expected bitset length 8, got %d", bs.Len())
		}
		if !bs.IsSet(7) {
			t.Error("expected bit 7 to be set")
		}
		if bs.IsSet(0) {
			t.Error("expected bit 0 to be unset")
		}
	})
}

func TestClone(t *testing.T) {
	bs := New(8)
	bs.Set(true, 1, 2)
	clone := Clone(bs)
	if !clone.IsSet(1) || !clone.IsSet(2) {
		t.Error("cloned bitset does not match original")
	}
	clone.Set(false, 1)
	if bs.IsSet(1) == clone.IsSet(1) {
		t.Error("expected original to remain unchanged")
	}
}

func TestParseBytes(t *testing.T) {
	_, err := ParseHexBytes([]byte("83"))
	if err != nil {
		t.Errorf("unexpected error parsing bytes: %v", err)
	}
	_, err = ParseHexBytes([]byte("8"))
	if err == nil {
		t.Error("expected error for odd-length byte slice")
	}
}

func TestValidate(t *testing.T) {

	t.Run("empty input", func(t *testing.T) {
		err := Validate([]byte{})
		if err != nil {
			t.Errorf("expected no error for empty input")
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		err := Validate([]byte("1z"))
		if err == nil {
			t.Error("expected error for invalid hex input")
		}
	})
	t.Run("odd-length input", func(t *testing.T) {
		err := Validate([]byte("1"))
		if err == nil {
			t.Error("expected error for odd-length input")
		}
	})

	t.Run("valid input", func(t *testing.T) {
		err := Validate([]byte("1a"))
		if err != nil {
			t.Errorf("expected valid input: %v", err)
		}

		err = Validate([]byte("1a2b"))
		if err != nil {
			t.Errorf("expected valid input: %v", err)
		}
	})
}

func TestAreSet(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		res, err := AreSet("b3", All)
		if err != nil {
			t.Error("expected no error, got:", err)
		}
		if res {
			t.Error("expected no bits to be checked")
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		res, err := AreSet("x1", All, 0, 1)
		if err == nil {
			t.Error("expected error for invalid hex string")
		}
		if res {
			t.Error("expected no bits to be set")
		}
	})

	t.Run("odd-length input", func(t *testing.T) {
		res, err := AreSet("b31", All, 0, 1)
		if err == nil {
			t.Error("expected error for odd-length hex string")
		}
		if res {
			t.Error("expected no bits to be set")
		}
	})
	t.Run("valid input: all", func(t *testing.T) {
		res, err := AreSet("b3", All, 0, 2, 3, 6, 7)
		if err != nil {
			t.Error("expected no error, got:", err)
		}
		if !res {
			t.Error("expected bits 0, 2, 3, 6, and 7 to be set")
		}
	})

	t.Run("valid input: all subset", func(t *testing.T) {

		res, err := AreSet("b3", All, 0, 6)
		if err != nil {
			t.Error("expected bit 6 to be unset")
			t.FailNow()
		}
		if !res {
			t.Error("expected bits 0 and 6 to be unset")
		}
	})
	t.Run("valid input: any", func(t *testing.T) {
		res, err := AreSet("b3", Any, 1, 4, 6)
		if err != nil {
			t.Error("expected no error, got:", err)
		}
		if !res {
			t.Error("expected bits 0, 1, or 2 to be set")
		}
	})
	t.Run("valid input: any, invalid subset", func(t *testing.T) {
		res, err := AreSet("b3", Any, 5, 15)
		if err != nil {
			t.Error("expected no error, got:", err)
		}
		if res {
			t.Error("expected bits 5 and 15 to be unset")
		}
	})
}

func TestByteBitSet_BinaryString(t *testing.T) {
	t.Run("empty bit set", func(t *testing.T) {
		bs := New(0)
		str := bs.BinaryString()
		if str != "" {
			t.Errorf("expected empty string, got %s", str)
		}
	})
	t.Run("valid bit set", func(t *testing.T) {
		bs := New(8)
		bs.Set(true, 0, 6)
		expected := "10000010"
		if str := bs.BinaryString(); str != expected {
			t.Errorf("expected %s, got %s", expected, str)
		}
	})
}

func TestByteBitSet_Set(t *testing.T) {

	t.Run("zero length, short mask", func(t *testing.T) {
		bs := New(0)
		bs.Set(true, 5)
		if !bs.IsSet(5) {
			t.Error("expected bit 5 to be set after auto-extension")
		}
	})

	t.Run("zero length, long mask", func(t *testing.T) {
		bs := New(0)
		bs.Set(true, 1000)
		if !bs.IsSet(1000) {
			t.Error("expected bit 1000 to be set after auto-extension")
		}
	})

	t.Run("out of bounds", func(t *testing.T) {
		bs := New(8)
		defer func() {
			if r := recover(); r != nil {
				t.Error("expected no panic for out of bounds set")
			}
		}()

		bs.Set(true, 100)
		if !bs.IsSet(100) {
			t.Error("expected bit 1000 to be set after auto-extension")
		}
	})
}

func TestParseString(t *testing.T) {

	t.Run("empty string", func(t *testing.T) {
		bs, err := ParseHexString("")
		if err != nil {
			t.Errorf("unexpected error for empty string: %v", err)
		}

		if bs.Len() != 0 {
			t.Error("expected zero-length bitset")
		}
	})
	t.Run("invalid hex string", func(t *testing.T) {
		_, err := ParseHexString("aabz")
		if err == nil {
			t.Error("expected error for invalid hex string")
		}
	})
	t.Run("odd-length hex string", func(t *testing.T) {
		_, err := ParseHexString("abc")
		if err == nil {
			t.Error("expected error for odd-length hex string")
		}
	})

	t.Run("invalid due to non-hexadecimal runes", func(t *testing.T) {
		_, err := ParseHexString("1с你好")
		if err == nil {
			t.Error("expected error for non-hexadecimal runes")
		}
	})

	t.Run("valid hex string", func(t *testing.T) {
		bs, err := ParseHexString("1A")
		if err != nil {
			t.Errorf("unexpected error for valid hex string: %v", err)
		}
		if bs.Len() != 8 {
			t.Errorf("expected bitset length 8, got %d", bs.Len())
		}
		if !bs.AreSet(All, 3, 4, 6) {
			t.Error("expected bits 3, 4, and 6 to be set")
		}
	})
}

func BenchmarkByteBitSet_IsSet(b *testing.B) {
	bs := New(1000)
	bs.Set(true, 500) // Set a single bit for testing

	b.ResetTimer()
	for b.Loop() {
		_ = bs.IsSet(500)
	}
}

func BenchmarkByteBitSet_AreSet_One(b *testing.B) {
	bs := New(1000)
	bs.Set(true, 500) // Set multiple bits for testing

	b.ResetTimer()
	for b.Loop() {
		_ = bs.AreSet(All, 500)
	}
}

func BenchmarkByteBitSet_AreSet_Five(b *testing.B) {
	bs := New(1000)
	bs.Set(true, 500) // Set multiple bits for testing

	b.ResetTimer()
	for b.Loop() {
		_ = bs.AreSet(All, 500, 501, 502, 503, 504)
	}
}

func BenchmarkAreSet(b *testing.B) {
	hexString := "b3aaccddff115678901212121212121212121212121212121212121212121212121212121212121212"
	b.ResetTimer()
	for b.Loop() {
		_, _ = AreSet(hexString, All, 2, 3)
	}
}
