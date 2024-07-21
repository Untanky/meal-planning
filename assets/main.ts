import 'htmx.org';
import './main.css';
import {
    BarController,
    BarElement,
    Chart,
    ChartDataset,
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

type NutritionData = {
    date: DateTime;
    calories: number;
    weight: number;
};

const element = document.querySelector('#nutrition-diagram-section') as HTMLElement | null;
const ctx = document.querySelector('canvas#nutrition-diagram') as HTMLCanvasElement | null;

if (element && ctx) {
    const rawData = element.dataset.nutrition;
    if (rawData) {
        const data = JSON.parse(rawData) as Array<{ date: string; calories: number; weight: number; }>;
        const nutrition = data.map(element => ({ ...element, date: DateTime.fromISO(element.date) }))

        console.log(data[0]);

        const labels = data.map((nutrition) => nutrition.date);
        const caloriesDataset: ChartDataset = {
            type: 'bar',
            label: 'Calories (kCal)',
            borderColor: 'rgb(186, 230, 253)',
            backgroundColor: 'rgba(125, 211, 252, 0.5)',
            yAxisID: 'caloriesAxis',
            data: data.map((nutrition) => nutrition.calories),
        };
        const weightDataset: ChartDataset = {
            type: 'line',
            label: 'Weight (kg)',
            borderColor: 'rgba(253, 230, 138)',
            backgroundColor: 'rgba(252, 211, 77)',
            yAxisID: 'weightAxis',
            data: data.map((nutrition) => nutrition.weight),
        };

        new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [weightDataset, caloriesDataset],
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
