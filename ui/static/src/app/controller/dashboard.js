define([], function () {

    function controller($scope, $http, $interval) {
        $scope.showGraph = false
        $scope.chartOptions = {
            legend: {
                show: false
            },
            lines: {
                show: true,
                lineWidth: 4
            },
            points: {
                show: true,
                radius: 4
            },
            grid: {
                show: false,
                hoverable: true
            },
            shadowSize: 0,
            highlightColor: 10
        }
        $scope.chartData = []

        $http.get("/package").then(
            function(response) {
                $scope.packageName = response.data
            },
            function(response){}
        );

        function updateGraph() {
            $scope.i = 0
            $http.get("/fetch").then(function(response) {
                angular.forEach(response.data, function(benchmark) {
                    $scope.chartData[$scope.i] = {
                        points: [],
                        current : benchmark.Points[benchmark.Points.length - 1][1] / 1000000,
                        name: benchmark.Name,
                        change: (benchmark.Points[benchmark.Points.length - 1][1] - benchmark.Points[benchmark.Points.length - 2][1]) / 1000000
                    }
                    $scope.chartData[$scope.i].points.push({
                        label: benchmark.Name,
                        data: benchmark.Points,
                        color: "#FFFFFF"
                    })
                $scope.i++
                });
                $scope.showGraph = true
            },
            function(response) {
            });
        }
        updateGraph()
        $interval(updateGraph, 5000)
    }

    controller.$inject=['$scope', '$http', '$interval'];

    return controller;
});
