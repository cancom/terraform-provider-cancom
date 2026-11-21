package dynamiccloud

import (
	"reflect"
	"sort"
	"testing"
)

func TestIsIamUserCrn(t *testing.T) {
	valid := []string{
		"crn:cancom::iam:user:alice.smith@example.com",
		"crn:123456:region1:iam:user:bob-1@sub.example.de",
	}
	for _, crn := range valid {
		if !IsIamUserCrn(crn) {
			t.Errorf("expected %q to be a valid IAM user CRN", crn)
		}
	}

	invalid := []string{
		"",
		"alice@example.com",
		"crn:cancom::iam:role:admin",
		"crn:cancom::iam:user:no-at-sign",
	}
	for _, crn := range invalid {
		if IsIamUserCrn(crn) {
			t.Errorf("expected %q to be rejected as IAM user CRN", crn)
		}
	}
}

func TestIsVPCProjectServiceUserName(t *testing.T) {
	if !isVPCProjectServiceUserName("svc_t1-proj1_1", "t1", "proj1") {
		t.Error("expected valid service user name to match")
	}
	if !isVPCProjectServiceUserName("svc_t1-proj1_9999", "t1", "proj1") {
		t.Error("expected upper-bound service user name to match")
	}

	invalid := []string{
		"svc_t1-proj1_0",
		"svc_t1-proj1_10000",
		"svc_t1-other_1",
		"user_t1-proj1_1",
		"svc_t1-proj1_abc",
		"svc_t1-proj1_1@example.com",
	}
	for _, name := range invalid {
		if isVPCProjectServiceUserName(name, "t1", "proj1") {
			t.Errorf("expected %q not to match service user name", name)
		}
	}
}

func TestGetHumanAndServiceUsers(t *testing.T) {
	mixed := []string{
		"crn:t1::iam:user:alice@example.com",
		"svc_t1-proj_42",
		"crn:t1::iam:user:bob@example.com",
	}

	humans, err := getHumanUsers("t1", "proj", mixed)
	if err != nil {
		t.Fatalf("getHumanUsers failed: %v", err)
	}
	wantHumans := []string{
		"crn:t1::iam:user:alice@example.com",
		"crn:t1::iam:user:bob@example.com",
	}
	if !reflect.DeepEqual(humans, wantHumans) {
		t.Fatalf("getHumanUsers = %v, want %v", humans, wantHumans)
	}

	svcs, err := getServiceUsers("t1", "proj", mixed)
	if err != nil {
		t.Fatalf("getServiceUsers failed: %v", err)
	}
	wantSvcs := []string{"svc_t1-proj_42"}
	if !reflect.DeepEqual(svcs, wantSvcs) {
		t.Fatalf("getServiceUsers = %v, want %v", svcs, wantSvcs)
	}

	if _, err := getHumanUsers("t1", "proj", []string{"invalid-entry"}); err == nil {
		t.Fatal("expected error for invalid entry in getHumanUsers")
	}
	if _, err := getServiceUsers("t1", "proj", []string{"invalid-entry"}); err == nil {
		t.Fatal("expected error for invalid entry in getServiceUsers")
	}
}

func TestUsersSetRoundTrip(t *testing.T) {
	if got := setToUsers(nil); len(got) != 0 {
		t.Fatalf("expected empty slice for nil input, got %v", got)
	}

	in := []string{"crn:t1::iam:user:b@example.com", "crn:t1::iam:user:a@example.com"}
	out := setToUsers(usersToSet(in))
	sort.Strings(in)
	sort.Strings(out)
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round-trip mismatch: got %v, want %v", out, in)
	}
}
