document.addEventListener("DOMContentLoaded", function () {
    var cy = window.cy = cytoscape({
        container: document.getElementById('cy'),
        boxSelectionEnabled: false,
        autounselectify: true,
        style: [
            {
                selector: 'node',
                css: {
                    'width': '200',
                    'height': '50',
                    'color': 'white',
                    'content': 'data(id)',
                    'text-valign': 'center',
                    'text-halign': 'center'
                }
            },
            {
                selector: ':parent',
                css: {
                    'color': 'black',
                    'text-valign': 'top',
                    'text-halign': 'center',
                }
            },
            {
                selector: 'edge',
                css: {
                    'width': 'data(width)',
                    'curve-style': 'bezier',
                    'target-arrow-shape': 'triangle',
                    'line-style': 'dashed',
                    'line-dash-pattern': [4, 2],
                    'line-dash-offset': -0,
                    'color': '#B71C1C',
                },
            },
            {
                selector: "edge[label]",
                css: {
                    "label": "data(label)",
                    "text-rotation": "autorotate",
                    "text-margin-x": "0px",
                    "text-margin-y": "0px",
                    "text-background-opacity": 1,
                    "text-background-color": "#ffffff",
                }
            },
        ],

        elements: fetch('api/graph').then(function (res) {
            return res.json();
        }),

        layout: {
            name: 'cose-bilkent',
            animate: false,
            idealEdgeLength: 280,
        },
    });

    let offset = 0;

    function draw() {
        cy.edges().animate({
            duration: 0,
            style: {
                'line-dash-offset': -offset,
                'line-dash-pattern': [4, 2]
            }
        });
    }

    function march() {
        offset++;
        if (offset > 16) {
            offset = 0;
        }
        requestAnimationFrame(draw);
        setTimeout(march, 50);
    }

    march();
});