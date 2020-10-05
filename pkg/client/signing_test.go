package client_test

// Signing does not work yet.
// func TestEchoWithSigning(t *testing.T) {
//   c, err := client.NewAnyClientFromEnvs(false, nil)
//   if err != nil {
//     t.Fatalf("could not create client: %v", err)
//	 }
//	 ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
//	 defer cancel()
//	 if err := client.ExecuteEcho(ctx, c); err != nil {
//     t.Fatalf("echo test failed: %v", err)
//	 }
// }
