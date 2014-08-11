$(function() {
    console.log("Start");

    var data = [ { x: 0, y: 40 }, { x: 1, y: 49 }, { x: 2, y: 17 }, { x: 3, y: 42 } ];

    var graph = new Rickshaw.Graph( {
            element: document.querySelector("#chart"),
            width: 780,
            height: 250,
            series: [ {
                    color: 'steelblue',
                    data: data
            } ]
    } );

    graph.render();
});
