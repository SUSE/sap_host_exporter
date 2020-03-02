package internal

import "net/http"

func Landing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
<html>
<head>
	<title>SAP Host Exporter</title>
</head>
<body>
	<h1>SAP Host Exporter</h1>
	<h2>Prometheus exporter for SAP systems</h2>
	<ul>
		<li><a href="metrics">Metrics</a></li>
		<li><a href="https://github.com/SUSE/sap_host_exporter" target="_blank">GitHub</a></li>
	</ul>
</body>
</html>
`))
}
