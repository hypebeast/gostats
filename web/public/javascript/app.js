$(function() {
    var godocData = [];

    var getGodocStats = function (callback) {
        $.getJSON('/api/godocstats', function (data) {
            godocData = data;
            callback();
        });
    };

    getGodocStats(function () {
        var dataPoints = [];
        $.each(godocData, function (index, value) {
            dataPoints.push({x: value.date, y: value.count});
        });

        var min = dataPoints[0].y;
        var max = dataPoints[dataPoints.length-1].y;
        var graph = new Rickshaw.Graph({
            element: document.querySelector("#chart"),
            width: 580,
            height: 250,
            min: min,
            max: max,
            series: [ {
                color: 'steelblue',
                data: dataPoints
            } ]
        });

        var x_axis = new Rickshaw.Graph.Axis.Time( { graph: graph } );
        var y_axis = new Rickshaw.Graph.Axis.Y({
            graph: graph,
            orientation: 'left',
            tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
            element: document.getElementById('y-axis'),
        });
        graph.render();
    });
});
