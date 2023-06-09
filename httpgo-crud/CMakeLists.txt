cmake_minimum_required(VERSION 2.6 FATAL_ERROR)

project(httpgo-crud NONE)

if(NOT CMAKE_Go_COMPILER)
  if(NOT $ENV{GO_COMPILER} STREQUAL "")
    get_filename_component(CMAKE_Go_COMPILER_INIT $ENV{GO_COMPILER} PROGRAM PROGRAM_ARGS CMAKE_Go_FLAGS_ENV_INIT)

    if(CMAKE_Go_FLAGS_ENV_INIT)
      set(CMAKE_Go_COMPILER_ARG1 "${CMAKE_Go_FLAGS_ENV_INIT}" CACHE STRING "First argument to Go compiler")
    endif()

    if(NOT EXISTS ${CMAKE_Go_COMPILER_INIT})
      message(SEND_ERROR "Could not find compiler set in environment variable GO_COMPILER:\n$ENV{GO_COMPILER}.")
    endif()

  endif()

  set(Go_BIN_PATH
    $ENV{GOPATH}
    $ENV{GOROOT}
    $ENV{GOROOT}/../bin
    $ENV{GO_COMPILER}
    /usr/bin
    /usr/local/bin
    )

  if(CMAKE_Go_COMPILER_INIT)
    set(CMAKE_Go_COMPILER ${CMAKE_Go_COMPILER_INIT} CACHE PATH "Go Compiler")
  else()
    find_program(CMAKE_Go_COMPILER
      NAMES go
      PATHS ${Go_BIN_PATH}
    )
    EXEC_PROGRAM(${CMAKE_Go_COMPILER} ARGS version OUTPUT_VARIABLE GOLANG_VERSION)
    STRING(REGEX MATCH "go[0-9]+.[0-9]+.[0-9]+[ /A-Za-z0-9]*" VERSION "${GOLANG_VERSION}")
    message("-- The Golang compiler identification is ${VERSION}")
    message("-- Check for working Golang compiler: ${CMAKE_Go_COMPILER}")
  endif()

endif()

add_custom_command(
  OUTPUT ${PROJECT_NAME}${CMAKE_SHARED_LIBRARY_SUFFIX}
  DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/${PROJECT_NAME}.go
  COMMAND ${CMAKE_Go_COMPILER} build -o ${PROJECT_NAME}${CMAKE_SHARED_LIBRARY_SUFFIX} -buildmode=plugin -tags go_tarantool_ssl_disable -tags go_tarantool_call_17 ${CMAKE_CURRENT_SOURCE_DIR}/${PROJECT_NAME}.go
)

add_custom_target(lib ALL DEPENDS ${PROJECT_NAME}${CMAKE_SHARED_LIBRARY_SUFFIX})

## Install ####################################################################
###############################################################################

if(NOT DEFINED TARANTOOL_INSTALL_LUADIR)
    set(TARANTOOL_INSTALL_LUADIR "${PROJECT_SOURCE_DIR}/.rocks/share/tarantool")
endif()

if(NOT DEFINED TARANTOOL_INSTALL_LIBDIR)
    set(TARANTOOL_INSTALL_LIBDIR "${PROJECT_SOURCE_DIR}/.rocks/lib/tarantool")
endif()

configure_file(${CMAKE_CURRENT_SOURCE_DIR}/init.lua.in ${TARANTOOL_INSTALL_LUADIR}/${PROJECT_NAME}/init.lua)

install(
  FILES ${CMAKE_CURRENT_BINARY_DIR}/${PROJECT_NAME}${CMAKE_SHARED_LIBRARY_SUFFIX}
  DESTINATION ${TARANTOOL_INSTALL_LIBDIR}/${PROJECT_NAME}
)
