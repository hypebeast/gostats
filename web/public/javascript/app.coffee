$ ->
    getGodocMetrics = (cb) -> 
        $.getJSON '/api/godocstats', (data) ->
            cb(data)

    getGithubMetrics = (cb) ->
        $.getJSON '/api/godocstats', (data) ->
            cb(data)

    createGodocSeries = (data) ->
        _.map data, (d) ->
            {x: d.date, y: d.count }

    getMetrics = (cb) ->
        getGodocMetrics cb

    getMetrics (data) ->
        series = []
        data = createGodocSeries(data)
        series.push({data: data, color: 'rgba(211, 84, 0, 0.50)', stroke: 'rgba(211, 84, 0, 1.00)'})

        min = data[0].y
        max = data[data.length-1].y

        graph = new Rickshaw.Graph
            series: series,
            height: 200,
            renderer: 'area',
            stroke: true,
            min: min,
            max: max,
            element: document.querySelector('#chart')

        x_axis = new Rickshaw.Graph.Axis.Time
            graph: graph

        y_axis = new Rickshaw.Graph.Axis.Y
            graph: graph,
            orientation: 'left',
            tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
            element: document.querySelector('#y-axis')

        graph.render()

        countPackagesLastDay = data[data.length - 1].y - data[data.length - 2].y
        weekIndex = 8
        if data.length < 8
            weekIndex = data.length
        countPackagesLastWeek = data[data.length - 1].y - data[data.length - weekIndex].y
        totalPackages = data[data.length-1].y
        
        $('#newPackagesLastDay').text(countPackagesLastDay)
        $('#newPackagesLastWeek').text(countPackagesLastWeek)
        $('#totalPackages').text(totalPackages)

        resize = () ->
            newWidth = $('#chart').width() - 40
            graph.configure
                width: newWidth

            graph.render()

        window.addEventListener('resize', resize)
        resize()
