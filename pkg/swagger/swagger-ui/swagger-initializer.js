window.onload = function() {
  fetch("http://localhost:8080/swagger/services")
    .then(response => response.json())
    .then(data => {
      window.ui = SwaggerUIBundle({
        urls: data, // Use fetched data
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    })
    .catch(error => console.error("Failed to load Swagger services:", error));
};
