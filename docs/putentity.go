package bign_news

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

/*

PutEntity inserta una entidad en la tabla MySQL especificada utilizando el mapa de datos proporcionado.

data: representa un mapa de columnas y valores que se insertarán en la tabla

*/

func (p *MySQLRowsProvider) PutEntity(table string, data map[string]interface{}) (interface{}, error) {

	if len(data) == 0 {

		return nil, errors.New("no data to insert")

	}

	// Crear slices para las columnas y los valores que se insertarán y rellenar las slices con los datos del mapa

	columns := make([]string, 0, len(data))

	values := make([]interface{}, 0, len(data))

	for column, value := range data {

		columns = append(columns, column)

		values = append(values, value)

	}

	// Consulta SQL de inserción

	query := fmt.Sprintf("INSERT INTO %s (`%s`) VALUES (%s)",

		table,

		strings.Join(columns, "`, `"),

		strings.Repeat("?, ", len(values)-1)+"?",
	)

	stmt, err := p.DB.Prepare(query)

	if err != nil {

		return nil, err

	}

	// Incluyo un tiempo de descanso aunque no parece solucionarlo

	time.Sleep(3 * time.Second)

	defer stmt.Close()

	res, err := stmt.Exec(values...)

	// Incluyo un clierre de la conexión después de ejecutar

	if err != nil {

		return nil, err

	}

	defer stmt.Close()

	return res.LastInsertId()

}
