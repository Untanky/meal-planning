import 'htmx.org';
import * as d3 from 'd3';
import './main.css';
import {ScaleTime, timeFormat} from "d3";

const container = document.querySelector('#nutrition-diagram');

if (container) {
    // Declare the chart dimensions and margins.
    const width = container.getBoundingClientRect().width;
    const height = 300;
    const marginTop = 24;
    const marginRight = 24;
    const marginBottom = 24;
    const marginLeft = 48;

    const nutrition = [
        { date: new Date(2024, 6, 12), weight: 94.1, nutrition: 0 },
        { date: new Date(2024, 6, 13), weight: 93.75, nutrition: 0 },
        { date: new Date(2024, 6, 14), weight: 92.43, nutrition: 0 },
        { date: new Date(2024, 6, 15), weight: 91.46, nutrition: 0 },
        { date: new Date(2024, 6, 16), weight: 92.21, nutrition: 0 },
        { date: new Date(2024, 6, 17), weight: 92, nutrition: 0 },
        { date: new Date(2024, 6, 18), weight: 92, nutrition: 0 },
    ]

    // Declare the x (horizontal position) scale.
    const x = d3.scaleTime()
        .domain(d3.extent(nutrition, (n): Date => n.date) as [Date, Date])
        .range([marginLeft, width - marginRight]);

    // Declare the y (vertical position) scale.
    const y = d3.scaleLinear()
        .domain(d3.extent(nutrition, (n) => n.weight) as [number, number])
        .range([height - marginBottom, marginTop]);

    const weightLine = d3.line<{ date: Date; weight: number; }>()
        .x(d => x(d.date))
        .y(d => y(d.weight));

    // Create the SVG container.
    const svg = d3.create("svg")
        .attr("width", width)
        .attr("height", height);

    // Add the x-axis.
    svg.append("g")
        .attr("transform", `translate(0,${height - marginBottom})`)
        .call(d3.axisBottom(x).ticks(d3.timeDay));

    // Add the y-axis.
    svg.append("g")
        .attr("transform", `translate(${marginLeft},0)`)
        .call(d3.axisLeft(y));

    // Append a path for the line.
    svg.append("path")
        .attr("fill", "none")
        .attr("stroke", "steelblue")
        .attr("stroke-width", 1.5)
        .attr("d", weightLine(nutrition));

    // Append the SVG element.
    const node = svg.node();
    if (node) {
        container.append(node);
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
console.info(closeDialogButton);
if (closeDialogButton) {
    closeDialogButton.addEventListener('click', closeUpdateDialog);
}
