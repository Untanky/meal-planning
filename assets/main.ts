import 'htmx.org';
import './main.css';
import {
    BarController,
    BarElement,
    Chart,
    LinearScale,
    LineController,
    LineElement,
    PointElement,
    TimeScale
} from "chart.js";
import {DateTime, Duration} from 'luxon';
import 'chartjs-adapter-luxon';

Chart.register(BarController, BarElement, LinearScale, LineController, LineElement, PointElement, TimeScale);

const ctx = document.querySelector('canvas#nutrition-diagram') as HTMLCanvasElement | null;

if (ctx) {
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: [DateTime.now(), DateTime.now().minus(Duration.fromObject({day: 1}))],
            datasets: [
                {
                    type: 'bar',
                    label: 'Calories (kCal)',
                    borderColor: 'rgb(186, 230, 253)',
                    backgroundColor: 'rgba(125, 211, 252, 0.5)',
                    yAxisID: 'caloriesAxis',
                    data: [
                        {x: DateTime.local().minus(Duration.fromObject({day: 6})).startOf('day'), y: 2097},
                        {x: DateTime.local().minus(Duration.fromObject({day: 5})).startOf('day'), y: 1996},
                        {x: DateTime.local().minus(Duration.fromObject({day: 4})).startOf('day'), y: 2405},
                        {x: DateTime.local().minus(Duration.fromObject({day: 3})).startOf('day'), y: 2169},
                        {x: DateTime.local().minus(Duration.fromObject({day: 2})).startOf('day'), y: 2369},
                        {x: DateTime.local().minus(Duration.fromObject({day: 1})).startOf('day'), y: 2000},
                        {x: DateTime.local().minus(Duration.fromObject({day: 0})).startOf('day'), y: 2000},
                    ],
                },
                {
                    label: 'Weight (kg)',
                    borderColor: 'rgba(253, 230, 138, 0.5)',
                    backgroundColor: 'rgba(252, 211, 77)',
                    yAxisID: 'weightAxis',
                    data: [
                        {x: DateTime.local().minus(Duration.fromObject({day: 6})).startOf('day'), y: 94.1},
                        {x: DateTime.local().minus(Duration.fromObject({day: 5})).startOf('day'), y: 93.75},
                        {x: DateTime.local().minus(Duration.fromObject({day: 4})).startOf('day'), y: 92.43},
                        {x: DateTime.local().minus(Duration.fromObject({day: 3})).startOf('day'), y: 91.46},
                        {x: DateTime.local().minus(Duration.fromObject({day: 2})).startOf('day'), y: 92.21},
                        {x: DateTime.local().minus(Duration.fromObject({day: 1})).startOf('day'), y: 92},
                        {x: DateTime.local().minus(Duration.fromObject({day: 0})).startOf('day'), y: 92},
                    ],
                },
            ],
        },
        options: {
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'day',
                    },
                    grid: {
                        display: false,
                    }
                },
                weightAxis: {
                    title: {
                        text: 'Weight (kg)',
                        display: true,
                    },
                    grid: {
                        color: 'rgba(180, 83, 9, 0.2)',
                    },
                    position: 'left',
                },
                caloriesAxis: {
                    type: 'linear',
                    title: {
                        text: 'Calories (kCal)',
                        display: true,
                    },
                    grid: {
                        color: 'rgba(3, 105, 161, 0.2)',
                    },
                    beginAtZero: false,
                    ticks: {
                        stepSize: 100,
                    },
                    position: 'right'
                },
            },
        },
    });
}

const updateMealPlanDialog: HTMLDialogElement | null = document.querySelector('dialog#update-meal-day');

const openUpdateDialog = (event: Event) => {
    if (updateMealPlanDialog) {
        updateMealPlanDialog.showModal();
    }
};

const closeUpdateDialog = () => {
    if (updateMealPlanDialog) {
        updateMealPlanDialog.close();
    }
};

document.querySelectorAll('button.update-plan-button').forEach((element) => {
    element.addEventListener('click', openUpdateDialog);
});

const closeDialogButton = document.querySelector('button#close-update-meal-day-dialog');
if (closeDialogButton) {
    closeDialogButton.addEventListener('click', closeUpdateDialog);
}
