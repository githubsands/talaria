package money

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	money "github.com/Comcast/golang-money"
	"github.com/Comcast/webpa-common/wrp"
	"github.com/ugorji/go/codec"
)

// TalariaContainer holds the necessary primitives for talaria to complete a money span chain
type MoneyContainer struct {
	serverDecorator func(ctx context.Context, hs *money.HTTPSpanner, next http.Handler, r *http.Request, c money.MoneyContainer) http.Handler
	contents        chan []byte
	start           chan struct{}
	stop            chan struct{}
}

// TalariaON returns an HTTPSpannerOption that triggers a httpspanner into a talaria state
func TalariaON() money.HTTPSpannerOptions {
	return func(hs *money.HTTPSpanner) {
		hs.tr1d1um = nil
		hs.scytale = nil
		hs.talaria = TalariaContainer{
			serverDecorator: talariaServerDecorator,
			contents:        make(chan []byte),
			start:           make(chan struct{}),
			stop:            make(chan struct{}),
		}
	}
}

// TalariaServerDecorator decorates a http.handler with money.
func (t *TalariaContainer) ServerDecorator(ctx context.Context, hs *money.HTTPSpanner, next http.Handler, r *http.Request, c money.MoneyContainer) http.Handler {
	return talariaServerDecorator(ctx, hs, r, next, t)
}

func talariaServerDecorator(ctx context.Context, next http.Handler, r *http.Request, c money.MoneyContainer) http.Handler {
	tracker, err := money.ExtractTracker(r)

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return next.ServeHTTP(response, request)
	}

	htTracker := tracker.HTTPTracker()
	request = money.InjectTracker(r, htTracker)

	next.ServeHTTP(response, request)
	if err != nil {
		next.ServeHTTP(response, request)
	}

	var (
		message   = new(wrp.Message)
		stream, _ = ioutil.ReadAll(req.Body)
		mh        = codec.MsgpackHandle
		h         = &mh
	)

	if err := codec.NewEncoderBytes(&stream, h).Encode(message); err != nil {
		return next.ServeHTTP(response, request)
	}

	message.Traces = tracker.GetMaps()

	if err := codec.NewDecoderBytes(stream, h).Decode(&message); err != nil {
		return next.ServeHTTP(response, request)
	}

	// inject the request into the responses body.  Assumes that the contents of this request
	// can be encoded into a wrp message.
	request.Body = nopCloser{bytes.NewBuffer(stream)}

	// signal the deviceListenerDecorate to start executing
	<-c.start

	// run the next handler in a go routine.
	// this essentially acts as a transactor due to the nature of the read and write pumps.

	go next.ServeHTTP(response, request)

	// wait for traces here that are returned in by the event listener in the case that and event's wrp message contains a trace
	for {
		select {
		case bytes := <-c.contents:
			// TODO:

			//  1. decode contents
			if err := codec.NewEncoderBytes(&stream, h).Encode(message); err != nil {
				return nil, err
			}

			tracker.UpdateMaps(message.Traces)

			_, _ := tracker.Finish()

			// copy the http request, and clear it so no duplicate responses are sent back.
			moneyResponse := response
			for k := range moneyResponse.Header {
				delete(m, k)
			}

			// Inject the tracker
			request = InjectTracker(response, request)
			response.Request = request

			moneyHandler.ServeHTTP(rw, r)

			// close the traces channel
			close(traces)
			break

		// traces never made it bac
		default:
			close(traces)
			break

		}
	}
}

// NopCloser is a wrapper struct that allows our buffer to fullfil a closer method.
type nopCloser struct {
	io.Reader
}

// Close is used to allow NopCloser to fullfil the Closer method of a ReadCloser so we can fullfil
// the required methods of a http.Body
func (nopCloser) Close() error { return nil }
