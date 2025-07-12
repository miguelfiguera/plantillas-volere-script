package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Requisitos struct {
	Funcionales   []Funcional   `json:"requisitosFuncionales"`
	NoFuncionales []NoFuncional `json:"requisitosNoFuncionales"`
}

// Prioridad representa la prioridad de un requisito junto con su justificación
type Prioridad struct {
	Nivel         string `json:"nivel"`
	Justificacion string `json:"justificacion"`
}

type Funcional struct {
	ID                  string    `json:"id"`
	Nombre              string    `json:"nombre"`
	Modulo              string    `json:"modulo"`
	CasoDeUsoId         string    `json:"casoDeUsoId"`
	Descripcion         string    `json:"descripcion"`
	Prioridad           Prioridad `json:"prioridad"`
	CriteriosAceptacion []string  `json:"criteriosDeAceptacion"`
}

type NoFuncional struct {
	ID          string `json:"id"`
	Categoria   string `json:"categoria"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

func main() {
	baseDir, _ := filepath.Abs("../")
	jsonPath := filepath.Join(baseDir, "JsonRequisitos", "requisitos.json")

	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatal("Error leyendo archivo JSON:", err)
	}

	var reqs Requisitos
	if err := json.Unmarshal(data, &reqs); err != nil {
		log.Fatal("Error parseando JSON:", err)
	}

	outDir := filepath.Join(baseDir, "plantillasVolere")
	os.RemoveAll(outDir) // Limpiar directorio existente
	os.MkdirAll(outDir, os.FileMode(0755))

	// Generar plantillas para requisitos funcionales
	for _, rf := range reqs.Funcionales {
		filename := filepath.Join(outDir, fmt.Sprintf("%s_%s.md", rf.ID, sanitize(rf.Nombre)))
		content := buildFunctionalTemplate(rf)
		if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
			log.Printf("Error escribiendo archivo %s: %v", filename, err)
		}
	}

	// Generar plantillas para requisitos no funcionales
	for _, rnf := range reqs.NoFuncionales {
		filename := filepath.Join(outDir, fmt.Sprintf("%s_%s.md", rnf.ID, sanitize(rnf.Nombre)))
		content := buildNonFunctionalTemplate(rnf)
		if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
			log.Printf("Error escribiendo archivo %s: %v", filename, err)
		}
	}

	fmt.Printf("✓ Generadas %d plantillas funcionales\n", len(reqs.Funcionales))
	fmt.Printf("✓ Generadas %d plantillas no funcionales\n", len(reqs.NoFuncionales))
	fmt.Printf("✓ Total: %d plantillas Markdown en %s\n", len(reqs.Funcionales)+len(reqs.NoFuncionales), outDir)
	fmt.Println("\nPara convertir a PDF:")
	fmt.Println("  cd plantillasVolere")
	fmt.Println("  for file in *.md; do pandoc \"$file\" -o \"${file%.md}.pdf\"; done")
}

func sanitize(s string) string {
	// Convertir a minúsculas y reemplazar caracteres especiales
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	clean := reg.ReplaceAllString(s, "")
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(clean)), " ", "_")
}

func buildFunctionalTemplate(rf Funcional) string {
	criteriosList := ""
	for _, criterio := range rf.CriteriosAceptacion {
		criteriosList += fmt.Sprintf("- %s\n", criterio)
	}

	// Construir representación legible de la prioridad
	prioridadDisplay := rf.Prioridad.Nivel
	if rf.Prioridad.Justificacion != "" {
		prioridadDisplay = fmt.Sprintf("%s - %s", rf.Prioridad.Nivel, rf.Prioridad.Justificacion)
	}

	template := fmt.Sprintf(`<style>
@page {
    size: A4;
    margin: 2cm;
}

body {
    font-family: Arial, sans-serif;
    line-height: 1.6;
    color: #333;
    margin: 0;
    padding: 20px;
    background: white;
}

.container {
    max-width: 800px;
    margin: 0 auto;
    padding: 40px;
    min-height: 80vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.header {
    text-align: center;
    border-bottom: 3px solid #2c3e50;
    padding-bottom: 20px;
    margin-bottom: 40px;
}

.volere-title {
    font-size: 28px;
    font-weight: bold;
    color: #2c3e50;
    margin-bottom: 10px;
}

.requisito-id {
    font-size: 20px;
    font-weight: bold;
    color: #e74c3c;
    background: #f8f9fa;
    padding: 8px 16px;
    border-radius: 4px;
    display: inline-block;
}

.section {
    margin: 25px 0;
    text-align: left;
}

.section-title {
    font-size: 18px;
    font-weight: bold;
    color: #2c3e50;
    margin-bottom: 10px;
    border-left: 4px solid #3498db;
    padding-left: 15px;
}

.section-content {
    font-size: 14px;
    padding: 10px 0;
    text-align: justify;
}

.metadata {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    margin: 20px 0;
}

.metadata-item {
    background: #f8f9fa;
    padding: 15px;
    border-radius: 6px;
    border-left: 4px solid #3498db;
}

.metadata-label {
    font-weight: bold;
    color: #2c3e50;
    margin-bottom: 5px;
}

.criterios {
    background: #f8f9fa;
    padding: 20px;
    border-radius: 6px;
    border-left: 4px solid #27ae60;
}

.criterios ul {
    margin: 10px 0;
    padding-left: 20px;
}

.criterios li {
    margin-bottom: 8px;
    text-align: justify;
}

.footer {
    text-align: center;
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid #bdc3c7;
    font-size: 12px;
    color: #7f8c8d;
}

@media print {
    body {
        print-color-adjust: exact;
    }
    .container {
        page-break-inside: avoid;
    }
}
</style>

<div class="container">

<div class="header">

# Plantilla Volere

## %s

</div>

## Requisito

**%s**

### Metadatos

| Campo | Valor |
|-------|-------|
| **Módulo/Área** | %s |
| **Prioridad** | %s |
%s

## Descripción

%s

## Criterios de Aceptación

%s

## Fuente

Dra. Génesis / Equipo PeluDog

## Razón/Valor de Negocio

Este requisito contribuye a la digitalización y eficiencia operativa del consultorio veterinario PeluDog, mejorando la experiencia del cliente y la gestión interna de procesos.

---

<div class="footer">

*Plantilla Volere - Proyecto Socio Tecnologico: Plataforma CRM para PeluDog*

</div>

</div>`, rf.ID, rf.Nombre, rf.Modulo, prioridadDisplay,
		getCasoDeUsoSection(rf.CasoDeUsoId), rf.Descripcion, criteriosList)

	return template
}

func buildNonFunctionalTemplate(rnf NoFuncional) string {
	template := fmt.Sprintf(`<style>
@page {
    size: A4;
    margin: 2cm;
}

body {
    font-family: Arial, sans-serif;
    line-height: 1.6;
    color: #333;
    margin: 0;
    padding: 20px;
    background: white;
}

.container {
    max-width: 800px;
    margin: 0 auto;
    padding: 40px;
    min-height: 80vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.header {
    text-align: center;
    border-bottom: 3px solid #8e44ad;
    padding-bottom: 20px;
    margin-bottom: 40px;
}

.volere-title {
    font-size: 28px;
    font-weight: bold;
    color: #8e44ad;
    margin-bottom: 10px;
}

.requisito-id {
    font-size: 20px;
    font-weight: bold;
    color: #e74c3c;
    background: #f8f9fa;
    padding: 8px 16px;
    border-radius: 4px;
    display: inline-block;
}

.section {
    margin: 25px 0;
    text-align: left;
}

.section-title {
    font-size: 18px;
    font-weight: bold;
    color: #8e44ad;
    margin-bottom: 10px;
    border-left: 4px solid #9b59b6;
    padding-left: 15px;
}

.section-content {
    font-size: 14px;
    padding: 10px 0;
    text-align: justify;
}

.metadata {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
    margin: 20px 0;
}

.metadata-item {
    background: #f8f9fa;
    padding: 15px;
    border-radius: 6px;
    border-left: 4px solid #9b59b6;
}

.metadata-label {
    font-weight: bold;
    color: #8e44ad;
    margin-bottom: 5px;
}

.footer {
    text-align: center;
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid #bdc3c7;
    font-size: 12px;
    color: #7f8c8d;
}

@media print {
    body {
        print-color-adjust: exact;
    }
    .container {
        page-break-inside: avoid;
    }
}
</style>

<div class="container">

<div class="header">

# Plantilla Volere - Requisito No Funcional

## %s

</div>

## Requisito

**%s**

### Metadatos

| Campo | Valor |
|-------|-------|
| **Categoría** | %s |

## Descripción

%s

## Fuente

Dra. Génesis / Equipo PeluDog

## Razón/Valor de Negocio

Este requisito no funcional asegura la calidad, rendimiento y confiabilidad del sistema CRM, contribuyendo a una experiencia de usuario óptima y operaciones estables del consultorio.

---

<div class="footer">

*Plantilla Volere - Proyecto Socio Tecnologico: Plataforma CRM para PeluDog*

</div>

</div>`, rnf.ID, rnf.Nombre, rnf.Categoria, rnf.Descripcion)

	return template
}

func getCasoDeUsoSection(casoDeUsoId string) string {
	if casoDeUsoId == "" {
		return ""
	}
	return fmt.Sprintf("| **Caso de Uso** | %s |\n", casoDeUsoId)
}
