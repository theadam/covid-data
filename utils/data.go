package utils

func IsOrganization(province string) bool {
    return province == "US Military" ||
        province == "Federal Bureau of Prisons" ||
        province == "Veteran Hospitals"
}

func CountyIsOrganization(county string) bool {
    return county == "Michigan Department of Corrections (MDOC)" ||
        county == "Federal Correctional Institution (FCI)"
}
