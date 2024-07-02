import 'htmx.org';
import './main.css';

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
