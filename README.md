# irm-scrap



# Google Engine

  Para utilizar el engine de busqueda de google, se debe crear un objeto GoogleEngine

   NewGoogleAPIEngine(GoogleEngineOptions)

      Donde GoogleEngineOptions
            tiene las siguientes propiedades
            	RowsProvider any         <- No hace nada, lo dejé a modo de compatibilidad
	            Headers      []Header    <- Sonlos headers que se pueden agregar a la llamada http, 
                                            por omisión solo se carga "User-Agent"



  Teniendo ya creado el objeto

