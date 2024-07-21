import {DateTime} from "luxon";
import {
    BarController,
    BarElement,
    Chart,
    type ChartDataset,
    LinearScale,
    LineController,
    LineElement, PointElement, TimeScale, Tooltip
} from 'chart.js';
import 'chartjs-adapter-luxon';

type NutritionView = {
    date: string;
    calories?: number;
    weight?: number;
}

type NutritionData = {
    date: DateTime;
    calories: number | null;
    weight: number | null;
};

Chart.register(BarController, BarElement, LinearScale, LineController, LineElement, PointElement, TimeScale, Tooltip);

const getNutritionData = (element: HTMLElement): NutritionData[] => {
    const rawData = element.dataset.nutrition;
    if (!rawData) {
        throw new Error('attribute "data-nutrition" must exist on the selected element');
    }

    const data = JSON.parse(rawData) as Array<{ date: string; calories: number; weight: number; }>;
    return data.map(element => ({
        ...element,
        date: DateTime.fromISO(element.date),
    }));
}

const createNutritionChart = (canvas: HTMLCanvasElement, data: NutritionData[]): Chart => {
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

    return new Chart(canvas, {
        type: 'line',
        data: {
            labels,
            datasets: [caloriesDataset, weightDataset],
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
};

const registerTriggers = (chart: Chart) => {
    document.body.addEventListener('updateNutritionData', (event: Event & { detail?: NutritionView }) => {
        if (!event.detail) {
            throw new Error('must include detail field');
        }

        const data = {
            date: DateTime.fromISO(event.detail.date),
            calories: event.detail.calories || null,
            weight: event.detail.weight || null,
        } satisfies NutritionData;

        const index = chart.data.labels?.findIndex((label) => data.date.equals(label as DateTime));

        if (index === undefined || index === null  || index === -1) {
            console.warn('could not find element');
            return;
        }
        chart.data.datasets[0].data[index] = data.calories;
        chart.data.datasets[1].data[index] = data.weight;

        chart.update();
    });
};

export const setupNutritionDiagram = (selector: string) => {
    const canvas = document.querySelector<HTMLCanvasElement>(selector);

    if (!canvas) {
        throw new Error('selector must exist in the document');
    }

    const data = getNutritionData(canvas);
    const chart = createNutritionChart(canvas, data);
    registerTriggers(chart);
}

setupNutritionDiagram('canvas#nutrition-diagram');
