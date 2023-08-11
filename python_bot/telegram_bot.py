
import logging
from telegram import Update, ParseMode
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters, CallbackContext
import requests

TOKEN = '6647242605:AAGk3bwBMi3h4SSqXg5YhwEISBuzUM1DfX8'

# Настройка логирования
logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                     level=logging.INFO)
logger = logging.getLogger(__name__)

# Функция для обработки команды /start
def start(update: Update, _: CallbackContext) -> None:
    update.message.reply_text('Привет! Я TaskBotGo, бот для управления задачами. '
                              'Используйте команду /add для добавления задачи и /view для просмотра списка задач.')

# Функция для обработки команды /add
def add_task(update: Update, _: CallbackContext) -> None:
    # Получаем описание задачи из сообщения пользователя
    task_description = update.message.text.replace('/add', '').strip()
    if task_description:
        # Отправляем описание задачи в Golang API сервер для добавления в Redis
        response = requests.post("http://golang_api:8080/add", json={"description": task_description})
        if response.status_code == 201:
            update.message.reply_text('Задача успешно добавлена!')
        else:
            update.message.reply_text('Не удалось добавить задачу.')
    else:
        update.message.reply_text('Пожалуйста, укажите описание задачи после команды /add.')

# Функция для обработки команды /view
def view_tasks(update: Update, _: CallbackContext) -> None:
    # Отправляем запрос в Golang API сервер для получения списка задач из Redis
    response = requests.get("http://golang_api:8080/view")
    if response.status_code == 200:
        tasks = response.json()
        if tasks:
            task_list = '\n'.join(tasks)
            update.message.reply_text(f'Список задач:\n{task_list}')
        else:
            update.message.reply_text('Список задач пуст.')
    else:
        update.message.reply_text('Не удалось получить список задач.')

def main():
    # Укажите токен вашего Telegram бота
    updater = Updater(TOKEN)

    # Получаем диспетчер и добавляем обработчики команд
    dp = updater.dispatcher
    dp.add_handler(CommandHandler("start", start))
    dp.add_handler(CommandHandler("add", add_task))
    dp.add_handler(CommandHandler("view", view_tasks))

    # Запускаем бота
    updater.start_polling()
    updater.idle()

if __name__ == '__main__':
    main()
