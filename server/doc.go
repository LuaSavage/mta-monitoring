// Package server queries MTA:SA game servers over the ASE (All-Seeing Eye) UDP protocol.
//
// Create a Server with NewServer, then call UpdateOnce to fetch name, player count,
// and other fields. Set Server.Timeout (seconds) to control UDP deadlines.
package server
