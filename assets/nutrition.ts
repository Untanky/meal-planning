import {DateTime} from "luxon";
import {
    BarController,
    BarElement,
    Chart,
    ChartData,
    ChartDataset,
    LinearScale,
    LineController,
    LineElement, PointElement, TimeScale, Tooltip
} from 'chart.js';
import 'chartjs-adapter-luxon';

type NutritionData = {
    date: DateTime;
    calories: number;
    weight: number;
};

Chart.register(BarController, BarElement, LinearScale, LineController, LineElement, PointElement, TimeScale, Tooltip);

const parseNutritionData = (rawData: string): NutritionData[] => {
    const data = JSON.parse(rawData) as Array<{ date: string; calories: number; weight: number; }>;
    return data.map(element => ({
        ...element,
        date: DateTime.fromISO(element.date),
    }));
};

const createNutritionChart = (canvas: HTMLCanvasElement, data: ChartData): Chart => {
    return new Chart(canvas, {
        type: 'line',
        data,
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
};

const renderNutritionDiagram = (canvas: HTMLCanvasElement): void => {
    const rawData = canvas.dataset.nutrition;
    if (!rawData) {
        throw new Error('attribute "data-nutrition" must be set on parentHTML');
    }

    const data = parseNutritionData(rawData);

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

    const chart = createNutritionChart(canvas, {
        labels,
        datasets: [caloriesDataset, weightDataset]
    });
}

const ctx = document.querySelector('canvas#nutrition-diagram') as HTMLCanvasElement | null;
if (ctx) {
    renderNutritionDiagram(ctx);
}
