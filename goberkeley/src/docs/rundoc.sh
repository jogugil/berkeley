#!/bin/bash

# Directorios
SRC_DIR="src"
PACKAGE_DIR="src/berkeley"
OUTPUT_DIR="docs"
MODULE_PATH="github.com/jogugil/berkeley/goberkeley" # Cambia esto al nombre de tu módulo en go.mod

# Crear el directorio de salida para los archivos HTML
mkdir -p $OUTPUT_DIR

# Iniciar el servidor de godoc en segundo plano
echo "Iniciando el servidor de documentación Go..."
godoc -http=:6060 &

# Esperar a que el servidor esté listo
sleep 2

# Descargar la documentación generada para el paquete Berkeley y el archivo main.go
echo "Descargando documentación..."

# La URL debe incluir el nombre del módulo como parte del path
wget -r -O "$OUTPUT_DIR/berkeley_package.html" "http://localhost:6060/pkg/goberkeley/berkeley/"
wget -r -O "$OUTPUT_DIR/index_pkg.html" "http://localhost:6060/pkg/goberkeley/"
# Detener el servidor de godoc
kill $!

# Crear el archivo index.html
echo "Generando el índice de documentación..."

cat <<EOL > "$OUTPUT_DIR/index.html"
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Documentación de Go</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        h1 {
            color: #4CAF50;
        }
        ul {
            list-style-type: none;
            padding-left: 0;
        }
        li {
            margin: 5px 0;
        }
        a {
            color: #4CAF50;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>

    <h1>Documentación de Go</h1>
    <p>Bienvenido a la documentación de tu proyecto Go. Aquí puedes acceder a los archivos generados:</p>

    <ul>
        <li><a href="berkeley_package.html" target="_blank">Documentación del paquete Berkeley</a></li>
        <li><a href="main_go.html" target="_blank">Documentación del archivo main.go</a></li>
    </ul>

</body>
</html>
EOL

echo "Documentación generada con éxito en el directorio $OUTPUT_DIR."

