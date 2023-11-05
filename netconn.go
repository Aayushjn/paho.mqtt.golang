/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 *    MAtt Brittan
 */

package mqtt

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/quic-go/quic-go"
)

//
// This just establishes the network connection; once established the type of connection should be irrelevant
//

// openConnection opens a network connection using the protocol indicated in the URL.
// Does not carry out any MQTT specific handshakes.
func openConnection(uri *url.URL, tlsc *tls.Config, timeout time.Duration, headers http.Header, websocketOptions *WebsocketOptions, dialer *net.Dialer) (quic.Stream, error) {
	quicConfig := &quic.Config{
		Allow0RTT: true,
	}

	conn, err := quic.DialAddrEarly(context.Background(), uri.Host, tlsc, quicConfig)
	if err != nil {
		return nil, err
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	return stream, err
}
