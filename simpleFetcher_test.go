package main

const domainA = "domainA.com"
const domainB = "domainB.com"
const locationA = "https://" + domainA + "/pageA"
const locationA1 = "https://" + domainA + "/pageB"
const locationB = "https://" + domainB + "/pageB"

/*
func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content := []byte("Hello!")

	httpmock.RegisterResponder("GET", locationA,
		httpmock.NewBytesResponder(200, content))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	resp, _ := f.getURLContent(locationA)
	body, _ := ioutil.ReadAll(resp.Body)

	if !(reflect.DeepEqual(body, content)) {
		t.Fail()
	}
}

func TestRedirect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content := []byte("Hello!")

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationA1))
	httpmock.RegisterResponder("GET", locationA1,
		httpmock.NewBytesResponder(200, content))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	resp, _ := f.getURLContent(locationA)
	body, _ := ioutil.ReadAll(resp.Body)
	if !(reflect.DeepEqual(body, content)) {
		t.Fail()
	}
}

func TestInfiniteRedirect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationA))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	_, err := f.getURLContent(locationA)

	if err != errorRedirection {
		t.Fail()
	}
}

func TestDifferentDomainRedirect(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationB))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	_, err := f.getURLContent(locationA)

	if err != errorDomain {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	_, err := f.getURLContent(locationA)

	if err != errorFetching {
		t.Fail()
	}
}
*/
