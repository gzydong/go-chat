package adapter

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go-chat/internal/pkg/ichat/socket/adapter/encoding"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/strutil"
)

func TestTcp_Server(t *testing.T) {
	listener, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9501))

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		go func() {
			conn, err := NewTcpAdapter(conn)
			if err != nil {
				return
			}

			for {
				data, err := conn.Read()
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(data))
			}
		}()
	}
}

type Authorize struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

var jsons = `["eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIxMDQ2In0.fEV7hEma623iYi_hiE1tU_CGXoMt0rT2R-UklU2k7us","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIxMDU4In0.2zM4ekNh83PKmWP_rYkw7ZSewJvmlCAOge15U5__OqQ","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIxMDcwIn0.EnW4p9X72avpb-UVpPUinOnrTjwQemCOj9VE-ROg9rA","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIxMTQ5In0.PPmHyJMs2rgf2KzBIpRQ5kEAe2BPfZkaVTircjIunRc","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIxODg3In0.9DflYfVfQcOpQhmcK0xIXemSx1XAyiJ55Kae6TlWd_A","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIxOTAzIn0.0tPOVllDOtw3o2qk9971KHj55rtVdMD0oKrySgwhrEU","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIyMDAxIn0.cvPwAOQVfKdPQMbybjbvsL49OwdNDRCFyLQnAmT5AtU","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIyMDE2In0.576NXw_UgevqO_BFjCcmpJREx_aRTWUrNLKEiK0mDro","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIyMDIyIn0.XDtWDL2bT6r9ExD_jJlLboHfloJEC5BWZf7GOROKPeY","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIyMDUzIn0.qw_-91jaoMHmGP5vY9ZB8IHn3E3h96gmtTC67rwtDXg","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIyMDU1In0.-PomJvhl8MqKbWCngYCnhxajb8L5_zhLBOlfd380qaE","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMDE5In0.TGuC-PnjnmV8UCg6_L5W1_QaYllamZQt_Bgjim1r9Co","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMDQ1In0.M2wht4gx-APo01voTI6MKqG4yrC-39SlMCRLdru6cqI","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMDU0In0.4x_Su2htYONHNCEf-bNmsgtDm91D3S5hpsINFKMAHRw","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMDYzIn0.4HpHnsUDuxuIPjwLG7GXOJJsjxOAM8UOdgq4df9G4Ic","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMDcyIn0.zEeDNKdmP1Jp2yDppoimvlAwaoT63GbVGAkuheSwJWI","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMDg0In0.BxyOUhf1OmRn2L45Lh1nnqmHJ5gbQbRLOXg7YZXGuv8","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiIzMTE3In0.EDKueGwUbx1hrg3WDxmjBEhIKlNn1RMg8HMaBJpUzQo","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDE4In0.IyemdPtETxvsGFecbDN1lb5A5nhWlxLqZ-h0zHNnQbQ","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDMxIn0.f4MJbBvpT_Zz8QA9DRYhDQ-z1EcR8eAsFHqVS8DqyOk","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDQyIn0.V5EBofHsbeB9iXTLq1Il9YWtAcGdFQe6j2odyaPG9TM","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDQ0In0.FpLqSFTOaMJULlu7WnnxMHDJ66kFBERrO1Ojlfeoimo","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDQ1In0.2kiik9GfZC3So7jjkGWM4pAdjU34Fd9IjmyDKIKB2ZQ","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDQ2In0.QQ05Gnn44dtDdhWJUjB0zlYRmN1meLfUlVNelu03nWU","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDQ3In0.gmPBIcO0wFtPuyS7c0icLoUIggSsJsOxSUDYWY8cxMU","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDQ4In0.e81FA6DP0U7suw6kGgDIboYUI074_7G-EN8VF3iDDIM","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDY2In0.OsWEs4un7InwKFILyQD-TIf4wnJosNNo0vvG0cZfQ1M","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDcwIn0.--KvEnt2Tt8gi-Rd1R1dNym5M38iGbCztj8pxhQITU4","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDcxIn0.6kB1AgY_lq3DSjMhbAkQX4rQmwjLvkU411YVLmkFs7w","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDgzIn0.VDf_RkDE_kZWW20HZjNFT2pDI2d4k6eOYyyRTipJv6s","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MDk5In0.hyeT-CrywVP2uneCs7UiLm-F_eQWtQ5bCbwseGbSVrE","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTAwIn0.q4aOC61blXwrSwR3HwwWsu0-76JoKAimGWr82wmjuSk","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTAxIn0.v9TDpgwWCfi1nHEaImJATifllQUrCB3j7cSIWWmmPqA","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTAyIn0.n-6RqW4VzhZmBfuazEk9KzxUaCPEagB1J1Wus8oLg-M","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTAzIn0.FPtK-pXxV03QRcoq0y-W_QeWvAB-bc388uZTvOmIWUQ","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTA3In0.at1AEbv2UW5hjUG1jOcitCeS30nKI2deEukD9Rnukd0","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTA4In0.PvNsZWC6h4FIsIYzUnWcRp-PZtkUMBdpTkOE5Rng7P4","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTA5In0.FQmoRRr5sDn1N700NQSWcsyQ1hZxeFUxR0LREB6aR-E","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTYyIn0.wnOTnVpTAfkYVcvCkZzhACnYpxWwAqA5GMO92rzB9mY","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MTkyIn0.tRo99Sn-13jqJ00iHIPlmASqc8RULqNFe3Ck09k-hgM","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjE0In0.Yjt5b4HbCN_LmoKxcVUhCqn3d4gnH2RDoh1wjuSOwNs","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjE2In0.ZF2wpK2srzotDBeK-8IE7zFuQ0y2i2KGMJgbVgHiyu4","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjE4In0.bG9ebZO8wN6becFmJb-tENBSCfEO5YkusWlhmK8pQds","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjI1In0.Zzwb-oh2AlrkWlYPHTxCWRmQ0ZZPr9ZYFWiOZYOE5T8","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjQwIn0.YheNiiuP8vL2IxMIPI_R8TyNE2PlJqWjQsbkgbwvsTY","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjU0In0.1RnlNPvWDZXXadFgg1qWi7quvSj6_EPTIsaullh7Uiw","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjYxIn0.huqTnvP3o9uHDkUd6Uyolw339COvTLegYfcUChFVh60","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjYyIn0.BawVdeVfl5w0cKNIwPA4UbdU4r8gS-Dft7KlxBAOJmE","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjY4In0.uFcY-MpG7KSML23k-whiCF_He-fjbiIgDq4D4f80o8Y","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0Mjc5In0.VV10EFJdU4eFK_xCmWKDwUDTAyHkQ3nieyf2HAm8hig","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjgzIn0.YkyWLMYkjWXlHQ3sRbjTXJQoET208K_6k2R4Dpix640","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MjkyIn0.qU_UiPKcxoASmNp5r7jKnl1qXDUNCZmtTNR2-8LuMPg","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MzYyIn0.6ZLtATLrCnlsfRKU2TPh6FkgLK9jyQOLoe0fuyTVF6M","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MzY0In0.Uw-viKFjpHH4mTy9iUO1nrWyNQ8dqW4fwFQQOVemT_8","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MzY1In0.lRRmIEQd6rzIrcfb9ZtLedGNOJIrUN3h-bUqp3P2S80","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0MzcwIn0.GmJOO1StLPXwScfVK7w63a2KxiUcjYXdvE8mWmfyOxk","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0Mzc0In0.1eiAASegxL5WLqKdTFUSydd0jjpdsr_BI81k8HW-o2Q","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTY4MDUyODMwNywiaWF0IjoxNjc5NjY0MzA3LCJqdGkiOiI0NDA5In0.kcstm_hrEl-xn4BmJCI72hIR2nGhZFl1T-DOBrnARFE"]`

// nolint
func TestTcp_Client_new(t1 *testing.T) {

	tokens := make([]string, 0)

	jsonutil.Decode(jsons, &tokens)

	for _, val := range tokens {
		ct(val)
	}

	time.Sleep(1 * time.Hour)
}

// nolint
func ct(token string) {
	go func() {
		conn, err := net.Dial("tcp", "106.14.177.175:9505")
		if err != nil {
			fmt.Println("dial failed, err", err)
			return
		}

		defer conn.Close()

		data, _ := encoding.NewEncode(jsonutil.Marshal(Authorize{
			Token:   token,
			Channel: "chat",
		}))
		_, _ = conn.Write(data)

		go func() {
			for index := 0; index < 500; index++ {
				msg := fmt.Sprintf(`{"event":"im.message.publish","content":{"receiver":{"talk_type":1,"receiver_id":2054},"type":1,"content":"测阿珂神经%d内科","mention":{"all":0,"uids":[]}}}`, index)

				data, _ := encoding.NewEncode([]byte(msg))
				conn.Write(data)

				time.Sleep(10 * time.Millisecond)
			}
		}()

		go func() {
			for {
				data, err := encoding.NewDecode(conn)

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(data))
			}
		}()

		time.Sleep(1 * time.Hour)
	}()
}

// nolint
func TestTcp_Client2(t1 *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:9505")
		if err != nil {
			fmt.Println("dial failed, err", err)
			return
		}

		defer conn.Close()

		data, _ := encoding.NewEncode(jsonutil.Marshal(Authorize{
			Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTcxMzUwNDY2MywiaWF0IjoxNjc3NTA0NjYzLCJqdGkiOiIyMDU0In0.qMybbKgqYlDR-5mP7VnoDIV8ex9Hg_tsXu8cTekX7-c",
			Channel: "chat",
		}))
		_, _ = conn.Write(data)

		go func() {
			for {
				msg := fmt.Sprintf(`{"msg_id":"%s","event":"event.talk.text.message","body":{"receiver":{"talk_type":1,"receiver_id":2055},"type":1,"content":"测阿珂神经%d内科","mention":{"all":0,"uids":[]}}}`, strutil.NewMsgId(), 999999)

				data, _ := encoding.NewEncode([]byte(msg))
				conn.Write(data)

				time.Sleep(20 * time.Second)
			}
		}()

		go func() {

			for {
				data, err := encoding.NewDecode(conn)

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(data))
			}
		}()

		time.Sleep(1 * time.Hour)
	}()

	time.Sleep(50 * time.Minute)
}
