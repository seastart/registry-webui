function padNumber(num: number): string {
    return num.toString().padStart(2, '0')
}

/**
 * 格式化时间
 * @param stamp 时间戳s
 * @returns 
 */
export function formatDate(stamp: number): string {
    const date = new Date(stamp * 1000);
    const formattedDate = `${date.getFullYear()}-${padNumber(date.getMonth() + 1)}-${padNumber(date.getDate())} ${padNumber(date.getHours())}:${padNumber(date.getMinutes())}`;
    return formattedDate
}

/**
 * 格式化大小
 * @param bytes Bytes
 * @param decimals 小数点后几位
 * @returns 
 */
export function formatSize(bytes: number, decimals = 2): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}
