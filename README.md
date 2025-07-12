# Generador de Plantillas Volere

Este directorio contiene el script **`generate_volere.go`**, encargado de tomar un archivo JSON con los requisitos de un proyecto y convertir cada requisito en una _plantilla Volere_ en formato Markdown (ideal para luego exportar a PDF con _Pandoc_).

---

## 1. Estructura de Datos Esperada

El script espera un archivo JSON con la siguiente estructura de alto nivel:

```json
{
  "requisitosFuncionales": [ RequisitoFuncional, ... ],
  "requisitosNoFuncionales": [ RequisitoNoFuncional, ... ]
}
```

### 1.1 Requisito Funcional

```json
{
  "id": "RF-XXX",
  "nombre": "Nombre del requisito",
  "modulo": "Módulo al que pertenece",
  "casoDeUsoId": "CU-XXX",          // Opcional
  "descripcion": "Descripción detallada",
  "prioridad": {
    "nivel": "Alta | Media | Baja",
    "justificacion": "Texto explicando la prioridad" // Puede ser vacío
  },
  "criteriosDeAceptacion": [ "Criterio 1", "Criterio 2", ... ]
}
```

### 1.2 Requisito No Funcional

```json
{
  "id": "RNF-XXX",
  "categoria": "Rendimiento | Usabilidad | Seguridad | ...",
  "nombre": "Nombre del requisito",
  "descripcion": "Descripción detallada"
  // (Opcional) Puedes añadir el bloque "prioridad" con la misma forma que arriba
}
```

> **Nota**: A partir de esta versión la propiedad `prioridad` pasó de ser un _string_ simple a un objeto con `nivel` y `justificacion`. Si sigues usando un JSON antiguo, actualízalo (o adapta el script) antes de ejecutar.

---

## 2. Ejecución Rápida

1. Sitúate en la raíz del proyecto y construye el binario:

```bash
cd Scripts
go build -o volere-gen generate_volere.go
```

2. Ejecuta el generador:

```bash
./volere-gen
```

Se crearán archivos Markdown en `../plantillasVolere`. Cada archivo lleva como nombre el identificador del requisito y una versión _slug_ de su título.

### Exportar a PDF

Dentro de `plantillasVolere` puedes convertir todos los `.md` a `.pdf` con:

```bash
cd plantillasVolere
for file in *.md; do pandoc "$file" -o "${file%.md}.pdf"; done
```

Necesitas tener _Pandoc_ instalado en tu sistema.

---

## 3. Personalización para Otros Proyectos

1. **Estructura JSON**: Respeta la estructura mostrada en la sección 1. Cambia los campos, añade los que necesites, o borra los que no apliquen. El script ignora propiedades desconocidas.
2. **Plantilla HTML/CSS**: Dentro de `generate_volere.go` hay dos funciones (`buildFunctionalTemplate` y `buildNonFunctionalTemplate`) que contienen la plantilla embebida. Puedes modificar los estilos, los encabezados o la información mostrada.
3. **Directorios de salida**: Si quieres generar los documentos en otra carpeta, cambia la variable `outDir` en `main()`.
4. **Campos Adicionales**: Si tu JSON incluye campos extras (por ejemplo, _Estado_, _Fecha_...), añade esos campos en las structs correspondientes y actualiza la plantilla para mostrarlos.

---

## 4. Dependencias

- Go ≥ 1.18 (se usa la librería estándar exclusivamente).
- Opcional: [Pandoc](https://pandoc.org/) para convertir los Markdown a PDF.

---

## 5. Personalización del Encabezado y Pie de Página

A partir de la versión **1.1** es posible **cambiar fácilmente el título del proyecto y el texto del pie de página** que aparecen en las plantillas generadas.

1. Abre `Scripts/generate_volere.go` y localiza las funciones `buildFunctionalTemplate` y `buildNonFunctionalTemplate`.
2. Dentro de cada función busca el bloque **`<div class="header">`**. Allí encontrarás:
   - Un encabezado de primer nivel `# Plantilla Volere` (puedes cambiarlo por el _nombre de tu proyecto_).
   - Un encabezado de segundo nivel `## %s` que el script reemplaza en ejecución por el **identificador** del requisito. Si deseas anteponer el nombre del proyecto puedes, por ejemplo, dejarlo así:
     ```html
     # Plataforma CRM – PeluDog ## %s
     ```
3. Para **modificar el pie de página**, desplázate hasta el final de la misma cadena _template_ y ubica el bloque **`<div class="footer">`**. Cambia el texto entre las etiquetas `<div class="footer"> ... </div>` por el mensaje que desees. Ej.:
   ```html
   *Plantilla Volere – Proyecto XYZ – Versión 0.9*
   ```
4. Guarda los cambios y vuelve a compilar el binario si es necesario (`go build -o volere-gen generate_volere.go`).

> Las modificaciones descritas aplican tanto para los requisitos **funcionales** como para los **no funcionales**, ya que ambos utilizan la misma estructura interna de encabezado y pie de página.

---

¡Disfruta creando especificaciones Volere automatizadas! 🎉
