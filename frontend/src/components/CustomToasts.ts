import Swal from 'sweetalert2';

export function showSuccessToast(message: string): void {
    Swal.fire({
        toast: true,
        position: 'top-end',
        icon: 'success',
        title: message,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        background: '#4CAF50',
        color: '#fff'
    });
}



export function showErrorToast(message: string): void {
    Swal.fire({
        toast: true,
        position: 'top-end',
        icon: 'error',
        title: message,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        background: '#f44336', // red color
        color: '#fff'
    });
}