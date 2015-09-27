define([], function () {

    function controller($scope, $http, $interval) {
        $scope.showGraph = false
        $scope.chartOptions = {
            legend: {
                show: false
            },
            lines: {
                show: true,
                lineWidth: 10
            },
            splines: {
                show: true
            },
            grid: {
                show: false,
                hoverable: true
            },
            shadowSize: 0,
            highlightColor: 10
        }
        function updateGraph() {
            $scope.chartData = []
            $http.get("/fetch").then(function(response) {
                angular.forEach(response.data, function(benchmark) {
                    console.log(benchmark)
                    $scope.chartData.push({
                        label: benchmark.Name,
                        data: benchmark.Points,
                        color: "#FFFFFF"
                    })
                });
                $scope.showGraph = true
            },
            function(response) {
            });
        }
        $interval(updateGraph, 5000)

    }

    controller.$inject=['$scope', '$http', '$interval'];

    return controller;
});
