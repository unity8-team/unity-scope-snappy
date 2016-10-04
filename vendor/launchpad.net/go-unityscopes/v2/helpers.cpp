#include "helpers.h"

#include <cstdlib>
#include <cstring>

extern "C" {
#include "_cgo_export.h"
}

namespace gounityscopes {
namespace internal {

std::string from_gostring(const StrData str) {
    return std::string(str.data, str.length);
}

std::vector<std::string> split_strings(const StrData str) {
    std::vector<std::string> list;
    const char *s = str.data;
    // str contains a sequence of nul-terminated strings concatenated together.
    for (const char *p = str.data; p != str.data + str.length; ++p) {
        if (*p == '\0') {
            list.push_back(s);
            s = p + 1;
        }
    }
    return list;
}

void *as_bytes(const std::string &str, int *length) {
    *length = str.size();
    void *data = malloc(str.size());
    if (data == nullptr) {
        return nullptr;
    }
    memcpy(data, str.data(), str.size());
    return data;
}

}
}
