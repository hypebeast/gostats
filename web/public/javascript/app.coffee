$ ->
    getShortDate = (timestamp) ->
        date = moment.unix(timestamp)
        return Date.UTC(date.year(), date.month(), date.date())

    getGodocMetrics = (cb) -> 
        $.getJSON '/api/godocstats', (data) ->
            cb(data)

    getGithubRepoStarsHistory = (name, cb) ->
        $.getJSON '/api/github/stars?name=' + name, (data) ->
            cb(data)

    createGodocSeries = (data) ->
        _.map data, (d) ->
            [getShortDate(d.date), d.count]

    createGodocDailySeries = (data) ->
        _.map data, (d, key) ->
            if key < data.length-1
                count = data[key + 1].count - d.count
                [getShortDate(d.date), count]

    createGodocWeeklySeries = (data) ->
        series = []
        dateFirstMonday = Date.now
        # Set last day to the current day
        lastDay = data.length-1
        for i in [data.length-1..0]
            date = moment.unix(data[i].date)
            # If we found a monday, calculate the number of packages of this week
            if date.day() == 1
                series.push(data[lastDay].count - data[i].count)
                lastDay = i - 1
                dateFirstModay = getShortDate(data[i].date)

        return {series: series.reverse(), startDate: dateFirstModay}

    createTotalStarsSeries = (data) ->
        _.map data, (d) ->
            [getShortDate(d.date), d.stars]

    createStarsDailySeries = (data) ->
        _.map data, (d, key) ->
            if key < data.length-1
                count = data[key + 1].stars - d.stars
                [getShortDate(d.date), count]

    createChartGodocTotal = (data) ->
        godocTotalPackages = createGodocSeries(data)
        seriesTotalPackages = []
        seriesTotalPackages.push
            data: godocTotalPackages
            type: 'area'
            name: 'Total Packages'

        $('#chart-total-packages').highcharts
            title:
                text: 'Total Packages'
            xAxis:
                type: 'datetime'
            yAxis:
                title:
                    text: 'Count'
                min: godocTotalPackages[0][1]
                labels:
                    formatter: -> return this.value / 1000 + 'k'
            plotOptions:
                series:
                    lineColor: 'rgba(211, 84, 0, 1.0)'
                    fillColor: 'rgba(211, 84, 0, 0.5)'
                    marker:
                        fillColor: 'rgba(211, 84, 0, 1.0)'
            legend:
                enabled: false
            series: seriesTotalPackages

    createChartGodocDaily = (data) ->
        godocDailyData = createGodocDailySeries(data)
        seriesDailyData = []
        seriesDailyData.push
            data: godocDailyData
            type: 'area'
            name: 'Total Packages'

        $('#chart-daily-packages').highcharts
            title:
                text: 'Daily New Packages'
            xAxis:
                type: 'datetime'
            yAxis:
                title:
                    text: 'Count'
            plotOptions:
                series:
                    lineColor: 'rgba(211, 84, 0, 1.0)'
                    fillColor: 'rgba(211, 84, 0, 0.5)'
                    marker:
                        fillColor: 'rgba(211, 84, 0, 1.0)'
            legend:
                enabled: false
            series: seriesDailyData

    createChartGodocWeekly = (data) ->
        godocWeeklyData = createGodocWeeklySeries(data)
        seriesWeeklyData = []
        seriesWeeklyData.push
            data: godocWeeklyData.series
            type: 'area'
            name: 'Packages'
            pointStart: godocWeeklyData.startDate
            pointInterval: 7 * 24 * 3600 * 1000

        $('#chart-weekly-packages').highcharts
            title:
                text: 'Weekly New Packages'
            xAxis:
                type: 'datetime'
            yAxis:
                title:
                    text: 'Count'
            plotOptions:
                series:
                    lineColor: 'rgba(211, 84, 0, 1.0)'
                    fillColor: 'rgba(211, 84, 0, 0.5)'
                    marker:
                        fillColor: 'rgba(211, 84, 0, 1.0)'
            legend:
                enabled: false
            series: seriesWeeklyData

    godocSummaryStats = (data) ->
        countPackagesLastDay = data[data.length - 1].count - data[data.length - 2].count
        weekIndex = 8
        if data.length < 8
            weekIndex = data.length
        countPackagesLastWeek = data[data.length - 1].count - data[data.length - weekIndex].count
        totalPackages = data[data.length-1].count
        
        $('#newPackagesLastDay').text(countPackagesLastDay)
        $('#newPackagesLastWeek').text(countPackagesLastWeek)
        $('#totalPackages').text(totalPackages)

    getMetrics = (cb) ->
        getGodocMetrics cb

    getMetrics (data) ->
        createChartGodocTotal data
        createChartGodocDaily data
        createChartGodocWeekly data

        godocSummaryStats data

    createRepoCharts = (repoName) ->
        console.log repoName

    $('.button-stats-modal').click () ->
        repoName = $(this).attr('data-title')

        getGithubRepoStarsHistory repoName, (data) ->
            if data? and data.data? and data.data.length > 0
                seriesData = createTotalStarsSeries(data.data).reverse()
                minVal = seriesData[0][1]
                series = []
                series.push
                    data: seriesData
                    type: 'area'
                    name: 'Stars'

                $('#chart-total-stars').highcharts
                    title:
                        text: 'Total Stars'
                    xAxis:
                        type: 'datetime'
                    yAxis:
                        title:
                            text: 'Count'
                        min: minVal
                        labels:
                            formatter: -> return this.value / 1000 + 'k'
                    plotOptions:
                        series:
                            lineColor: 'rgba(211, 84, 0, 1.0)'
                            fillColor: 'rgba(211, 84, 0, 0.5)'
                            marker:
                                fillColor: 'rgba(211, 84, 0, 1.0)'
                    legend:
                        enabled: false
                    reflow: true
                    series: series

                # Calculate the number of new packages on a per day basis
                dailySeriesData = createStarsDailySeries(data.data.reverse())
                series = []
                series.push
                    data: dailySeriesData
                    type: 'column'
                    name: 'Count'

                $('#chart-daily-stars').highcharts
                    title:
                        text: 'Daily New Stars'
                    xAxis:
                        type: 'datetime'
                    yAxis:
                        title:
                            text: 'Count'
                    legend:
                        enabled: false
                    series: series
            else
                $('#chart-total-stars').text 'No data available'

            # Set title for modal and show it
            $('#modalTitle').text(repoName)
            $('#statsModal').modal 'show'

            reflowCharts = () ->
                $('#chart-total-stars').highcharts().reflow()
                $('#chart-daily-stars').highcharts().reflow()

            setTimeout ( ->
                reflowCharts()
            ), 200

