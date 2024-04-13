#!/bin/bash

# Endpoint URL
BASE_URL="http://localhost:8080"
POST_GET_SHORT_URL="/"
POST_GET_JSON_SHORT_URL="/api/shorten"
POST_GET_BATCH_SHORT_URL="/api/shorten/batch"
GET_USER_URLS="/api/user/urls"

function get_cookies(){
# Отправка запроса и извлечение куки из заголовка Set-Cookie
response=$(curl -i -X POST -H "Cookie: $COOKIES" "$BASE_URL$POST_GET_SHORT_URL")
local COOKIES=$(echo "$response" | grep -i 'Set-Cookie' | awk '{print $2}')
echo "$COOKIES"
}



# Function to make asynchronous requests
function make_async_requests_GetShortUrl() {
  local count="$1"
  COOKIE=$(get_cookies)
  for ((i = 1; i <= count; i++)); do
    # Generate a unique URL for each request (e.g., appending the iteration number to the base URL)
    current_url="$BASE_URL$POST_GET_SHORT_URL/$i"
	
   # Perform the request
    curl -X POST -H "Content-Type: text/plain" -H "Cookie: $COOKIES" -d "$current_url" "$BASE_URL$POST_GET_SHORT_URL" &
  done

  # Wait for all background jobs to finish
  wait
}

# Function to make asynchronous requests
function make_async_requests_GetJSONShortURL() {
  local count="$1"
  COOKIES=$(get_cookies)

  for ((i = 1; i <= count; i++)); do
    # Generate a unique URL for each request (e.g., appending the iteration number to the base URL)
    current_url="$BASE_URL$POST_GET_JSON_SHORT_URL/$i"
  
    # Perform the request
    curl -X POST -H "Content-Type: application/json" -H "Cookie: $COOKIE" -d '{"url":"'"$current_url"'"}' "$BASE_URL$POST_GET_JSON_SHORT_URL" &
  done

  # Wait for all background jobs to finish
  wait
}

# Function to make asynchronous requests
function make_async_requests_GetBatchShortURL() {
  
  local count="$1"

  for ((i = 1; i <= count; i++)); do
  
    # Отправка запроса и извлечение куки из заголовка Set-Cookie в каждой итерации цикла, чтобы каждый запрос был от нового пользователя
	COOKIES=$(get_cookies)

    # Generate a unique URL for each request (e.g., appending the iteration number to the base URL)
    current_url="$BASE_URL$POST_GET_BATCH_SHORT_URL/$i"
	# Perform the request
	correlation_id=$((i))
    data='[{"correlation_id": "'"$correlation_id"'","original_url":"'$current_url'"},{"correlation_id": "'"$correlation_id+1"'","original_url":"'$current_url+1'"},{"correlation_id": "'"$correlation_id+2"'","original_url":"'$current_url+2'"},{"correlation_id": "'"$correlation_id"'","original_url":"'$current_url'"}]'
    curl -X POST -H "Content-Type: application/json"  -H "Cookie: $COOKIES" -d "$data" "$BASE_URL$POST_GET_BATCH_SHORT_URL" &
  done

  # Wait for all background jobs to finish
  wait
}
# Function to make asynchronous requests
function make_async_requests_GetUserURLS() {
  
  local count="$1"
  COOKIES=$(get_cookies)
  
  for ((i = 1; i <= count; i++)); do
   
    curl -X GET -H "Cookie: $COOKIES"  "$BASE_URL$GET_USER_URLS" &
  done

  # Wait for all background jobs to finish
  wait
}


# Load testing
make_async_requests_GetShortUrl 100 & # Запускаем первую функцию в фоновом режиме
make_async_requests_GetJSONShortURL 100 & # Запускаем вторую функцию в фоновом режиме
make_async_requests_GetBatchShortURL 100 & # Запускаем третью функцию в фоновом режиме
make_async_requests_GetUserURLS 100 & # Запускаем четвертую функцию в фоновом режиме

# Ждем завершения всех фоновых задач
wait