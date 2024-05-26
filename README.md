# irm-scrap



# Google Engine

  Para utilizar el engine de busqueda de google, se debe crear un objeto GoogleEngine

   NewGoogleAPIEngine(GoogleEngineOptions)

      Donde GoogleEngineOptions
            tiene las siguientes propiedades
            	RowsProvider any         <- No hace nada, lo dejé a modo de compatibilidad
	            Headers      []Header    <- Sonlos headers que se pueden agregar a la llamada http, 
                                            por omisión solo se carga "User-Agent"



  Teniendo ya creado el objeto, se puede llamar a buscar las URL con:

        	miGoogleEngine.Search(search, 1)

            Donde search es un objeto IRMExtraSearchRequest
                  Se utilizaran las siguientes propiedades.
                     Query: "ليونل مسي",       <- Texto a buscar (puede ser cualquier idioma)
                     DateFrom:                 <- Opcional: desde fecha
                     DateTo:                   <- Opcional: hasta fecha
                     Markets:  []string{"AE"}  <- Opcional: un valor estándar ISO 3166-1 alpha-2.
                     MaxUrls:                  <- Opcional: maxima url de resultasdos, omision 10
                 

