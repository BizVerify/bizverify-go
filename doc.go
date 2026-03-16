// Package bizverify provides a Go client for the BizVerify business entity verification API.
//
// Create a client with an API key:
//
//	client := bizverify.New(bizverify.WithAPIKey("bv_live_..."))
//
// Verify a business entity:
//
//	resp, err := client.Verification.Verify(ctx, bizverify.VerifyParams{
//	    EntityName:   "Acme Inc",
//	    Jurisdiction: "us-fl",
//	})
//
// Search for entities:
//
//	resp, err := client.Search.Find(ctx, bizverify.SearchParams{
//	    EntityName: "Acme",
//	})
package bizverify
