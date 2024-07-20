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
    TimeScale,
    Tooltip
} from "chart.js";
import {DateTime, Duration} from 'luxon';
import 'chartjs-adapter-luxon';

Chart.register(BarController, BarElement, LinearScale, LineController, LineElement, PointElement, TimeScale, Tooltip);

const ctx = document.querySelector('canvas#nutrition-diagram') as HTMLCanvasElement | null;

if (ctx) {
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: [
                DateTime.local().minus(Duration.fromObject({day: 6})).startOf('day'),
                DateTime.local().minus(Duration.fromObject({day: 5})).startOf('day'),
                DateTime.local().minus(Duration.fromObject({day: 4})).startOf('day'),
                DateTime.local().minus(Duration.fromObject({day: 3})).startOf('day'),
                DateTime.local().minus(Duration.fromObject({day: 2})).startOf('day'),
                DateTime.local().minus(Duration.fromObject({day: 1})).startOf('day'),
                DateTime.local().minus(Duration.fromObject({day: 0})).startOf('day'),
            ],
            datasets: [
                {
                    type: 'bar',
                    label: 'Calories (kCal)',
                    borderColor: 'rgb(186, 230, 253)',
                    backgroundColor: 'rgba(125, 211, 252, 0.5)',
                    yAxisID: 'caloriesAxis',
                    data: [2097, 1996, 2405, 2169, 2369, 2000, 2000],
                },
                {
                    label: 'Weight (kg)',
                    borderColor: 'rgba(253, 230, 138)',
                    backgroundColor: 'rgba(252, 211, 77)',
                    yAxisID: 'weightAxis',
                    data: [94.1, 93.75, 92.43, 91.46, 92.21, 92, 92],
                },
            ],
        },
        options: {
            interaction: {
                mode: 'index',
                intersect: false,
            },
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
